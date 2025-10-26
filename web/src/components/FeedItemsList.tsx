import type { FeedItem } from '../types/feed';
import { sanitizeHtml } from '../utils/sanitizeHtml';

type FeedItemsListProps = {
  items: FeedItem[];
  loading: boolean;
};

const formatPublished = (value: string) => {
  if (!value) {
    return 'Data desconhecida';
  }

  try {
    return new Date(value).toLocaleString('pt-BR');
  } catch {
    return 'Data desconhecida';
  }
};

const extractPlainText = (html: string) =>
  html.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();

const LoadingSkeleton = () => (
  <li className="feed-card feed-card--skeleton">
    <div className="skeleton skeleton__date" />
    <div className="skeleton skeleton__title" />
    <div className="skeleton skeleton__line" />
    <div className="skeleton skeleton__line skeleton__line--short" />
  </li>
);

export const FeedItemsList = ({ items, loading }: FeedItemsListProps) => {
  if (loading) {
    return (
      <ul className="feed-list" aria-label="Carregando itens do feed">
        <LoadingSkeleton />
        <LoadingSkeleton />
        <LoadingSkeleton />
      </ul>
    );
  }

  if (items.length === 0) {
    return <p className="muted-text">Nenhum item encontrado neste feed.</p>;
  }

  return (
    <ul className="feed-list">
      {items.map((item) => {
        const sanitized = sanitizeHtml(item.description || '');
        const fallbackText = extractPlainText(sanitized);
        const hasHtml = sanitized && sanitized !== fallbackText;

        return (
          <li className="feed-card" key={`${item.link}-${item.title}`}>
            <time dateTime={item.publishedAt}>{formatPublished(item.publishedAt)}</time>
            <h3>
              <a href={item.link} target="_blank" rel="noreferrer">
                {item.title || 'Item sem título'}
              </a>
            </h3>

            {hasHtml ? (
              <div
                className="feed-card__content"
                dangerouslySetInnerHTML={{ __html: sanitized }}
              />
            ) : (
              fallbackText && <p>{fallbackText}</p>
            )}
          </li>
        );
      })}
    </ul>
  );
};
