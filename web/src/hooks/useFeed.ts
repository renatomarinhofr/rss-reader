import { useCallback, useState } from 'react';

import type { FeedResponse } from '../types/feed';

type FetchResult = {
  data: FeedResponse | null;
  fetchedAt: string | null;
};

export const useFeed = () => {
  const [feed, setFeed] = useState<FeedResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdatedAt, setLastUpdatedAt] = useState<string | null>(null);

  const fetchFeed = useCallback(async (url: string): Promise<FetchResult> => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/feed?url=${encodeURIComponent(url)}`);
      const payload = await response.json();

      if (!response.ok) {
        const message = typeof payload?.error === 'string' ? payload.error : 'Erro ao carregar o feed';
        throw new Error(message);
      }

      const data = payload as FeedResponse;
      const timestamp = data.fetchedAt || new Date().toISOString();

      setFeed(data);
      setLastUpdatedAt(timestamp);

      return { data, fetchedAt: timestamp };
    } catch (err) {
      setFeed(null);
      setLastUpdatedAt(null);

      const message = err instanceof Error ? err.message : 'Erro inesperado ao carregar o feed';
      setError(message);

      return { data: null, fetchedAt: null };
    } finally {
      setLoading(false);
    }
  }, []);

  const resetError = useCallback(() => setError(null), []);

  return {
    feed,
    loading,
    error,
    lastUpdatedAt,
    fetchFeed,
    resetError,
  };
};
