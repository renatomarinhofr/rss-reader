import type { RecentFeed } from '../types/feed';

type RecentFeedsProps = {
  feeds: RecentFeed[];
  onSelect: (url: string) => void;
  onClear: () => void | Promise<void>;
};

const formatRelativeTime = (value: string) => {
  try {
    const formatter = new Intl.RelativeTimeFormat('pt-BR', { numeric: 'auto' });
    const diff = Date.now() - new Date(value).getTime();
    const minutes = Math.round(diff / (60 * 1000));
    if (Math.abs(minutes) < 60) {
      return formatter.format(-minutes, 'minute');
    }
    const hours = Math.round(diff / (60 * 60 * 1000));
    if (Math.abs(hours) < 24) {
      return formatter.format(-hours, 'hour');
    }
    const days = Math.round(diff / (24 * 60 * 60 * 1000));
    return formatter.format(-days, 'day');
  } catch {
    return '';
  }
};

export const RecentFeeds = ({ feeds, onSelect, onClear }: RecentFeedsProps) => {
  if (feeds.length === 0) {
    return null;
  }

  return (
    <section className="recent-feeds">
      <div className="recent-feeds__header">
        <h2>Feeds recentes</h2>
        <button type="button" className="link-button" onClick={onClear}>
          Limpar histórico
        </button>
      </div>
      <ul>
        {feeds.map((item) => (
          <li key={item.sourceUrl}>
            <button type="button" onClick={() => onSelect(item.sourceUrl)}>
              <span className="recent-feeds__title">{item.title || item.sourceUrl}</span>
              <span className="recent-feeds__meta">
                {item.sourceUrl}
                {item.fetchedAt && <small> - {formatRelativeTime(item.fetchedAt)}</small>}
              </span>
            </button>
          </li>
        ))}
      </ul>
    </section>
  );
};
