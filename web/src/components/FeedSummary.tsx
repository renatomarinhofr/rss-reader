import type { FeedResponse } from '../types/feed';

type FeedSummaryProps = {
  feed: FeedResponse;
  lastUpdatedAt: string | null;
};

const formatDateTime = (value: string | null) => {
  if (!value) {
    return null;
  }

  try {
    return new Intl.DateTimeFormat('pt-BR', {
      dateStyle: 'long',
      timeStyle: 'short',
    }).format(new Date(value));
  } catch {
    return null;
  }
};

export const FeedSummary = ({ feed, lastUpdatedAt }: FeedSummaryProps) => {
  const formattedUpdatedAt = formatDateTime(lastUpdatedAt);

  return (
    <header className="feed-summary">
      <div>
        <h2>{feed.title || 'Feed sem título'}</h2>
        {feed.description && <p className="feed-summary__description">{feed.description}</p>}
        {feed.link && (
          <a className="feed-summary__link" href={feed.link} target="_blank" rel="noreferrer">
            Visitar website
          </a>
        )}
      </div>

      {formattedUpdatedAt && (
        <dl className="feed-summary__meta">
          <div>
            <dt>Atualizado</dt>
            <dd>{formattedUpdatedAt}</dd>
          </div>
          <div>
            <dt>Total de itens</dt>
            <dd>{feed.items.length}</dd>
          </div>
        </dl>
      )}
    </header>
  );
};
