import { useCallback, useEffect, useState } from 'react';

import { FeedCatalog } from './components/FeedCatalog';
import { FeedForm } from './components/FeedForm';
import { FeedItemsList } from './components/FeedItemsList';
import { FeedSummary } from './components/FeedSummary';
import { RecentFeeds } from './components/RecentFeeds';
import { StatusBanner } from './components/StatusBanner';
import { FEED_CATEGORIES, type FeedSource } from './data/sources';
import { useFeed } from './hooks/useFeed';
import { useRecentFeeds } from './hooks/useRecentFeeds';

const DEFAULT_CATEGORY = FEED_CATEGORIES[0];
const DEFAULT_FEED = DEFAULT_CATEGORY?.feeds[0];
const FALLBACK_URL = 'https://g1.globo.com/rss/g1/';
const INITIAL_URL = DEFAULT_FEED?.url ?? FALLBACK_URL;

const App = () => {
  const [url, setUrl] = useState(INITIAL_URL);
  const [activeCategoryId, setActiveCategoryId] = useState(
    DEFAULT_CATEGORY?.id ?? FEED_CATEGORIES[0]?.id ?? 'custom',
  );
  const [selectedCatalogFeed, setSelectedCatalogFeed] = useState<FeedSource | null>(
    DEFAULT_FEED ?? null,
  );
  const { feed, loading, error, lastUpdatedAt, fetchFeed, resetError } = useFeed();
  const { recentFeeds, registerRecentFeed, clearRecentFeeds } = useRecentFeeds();

  const alignCatalogSelection = useCallback((targetUrl: string) => {
    const trimmed = targetUrl.trim();
    if (!trimmed) {
      setSelectedCatalogFeed(null);
      return;
    }

    for (const category of FEED_CATEGORIES) {
      const found = category.feeds.find((feedSource) => feedSource.url === trimmed);
      if (found) {
        setActiveCategoryId(category.id);
        setSelectedCatalogFeed(found);
        return;
      }
    }

    setSelectedCatalogFeed(null);
  }, []);

  const handleLoadFeed = useCallback(
    async (targetUrl: string) => {
      if (!targetUrl) {
        return;
      }

      alignCatalogSelection(targetUrl);
      const { data, fetchedAt } = await fetchFeed(targetUrl);
      if (data && fetchedAt) {
        registerRecentFeed({
          sourceUrl: data.sourceUrl || targetUrl,
          title: data.title || targetUrl,
          description: data.description,
          link: data.link,
          fetchedAt,
        });
      }
    },
    [alignCatalogSelection, fetchFeed, registerRecentFeed],
  );

  useEffect(() => {
    void handleLoadFeed(INITIAL_URL);
  }, [handleLoadFeed]);

  const handleInputChange = (value: string) => {
    if (error) {
      resetError();
    }
    setUrl(value);
    setSelectedCatalogFeed(null);
  };

  const handleSubmit = (value: string) => {
    void handleLoadFeed(value);
  };

  const handleSelectRecent = (selectedUrl: string) => {
    setUrl(selectedUrl);
    void handleLoadFeed(selectedUrl);
  };

  const handleCategorySelect = (categoryId: string) => {
    setActiveCategoryId(categoryId);
  };

  const handleFeedSelect = (categoryId: string, feedSource: FeedSource) => {
    setActiveCategoryId(categoryId);
    setSelectedCatalogFeed(feedSource);
    setUrl(feedSource.url);
    void handleLoadFeed(feedSource.url);
  };

  return (
    <div className="app-container">
      <header className="page-header">
        <h1>Leitor de RSS</h1>
        <p>Explore fontes brasileiras confiáveis, salve seus favoritos e acompanhe as últimas matérias.</p>
      </header>

      <FeedCatalog
        categories={FEED_CATEGORIES}
        activeCategoryId={activeCategoryId}
        selectedFeedUrl={selectedCatalogFeed?.url}
        onCategorySelect={handleCategorySelect}
        onFeedSelect={handleFeedSelect}
      />

      <FeedForm value={url} loading={loading} onChange={handleInputChange} onSubmit={handleSubmit} />

      {error && <StatusBanner message={error} tone="error" onDismiss={resetError} />}

      <RecentFeeds feeds={recentFeeds} onSelect={handleSelectRecent} onClear={clearRecentFeeds} />

      {feed ? (
        <>
          <FeedSummary feed={feed} lastUpdatedAt={lastUpdatedAt} />
          <FeedItemsList items={feed.items} loading={loading} />
        </>
      ) : (
        !loading && <p className="muted-text">Informe uma URL válida para carregar um feed.</p>
      )}
    </div>
  );
};

export default App;
