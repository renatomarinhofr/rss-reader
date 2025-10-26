import { useCallback, useEffect, useState } from 'react';

import type { RecentFeed } from '../types/feed';

const MAX_ITEMS = 10;

type RecentFeedsResponse = {
  feeds: RecentFeed[];
};

export const useRecentFeeds = () => {
  const [recentFeeds, setRecentFeeds] = useState<RecentFeed[]>([]);

  const loadFeeds = useCallback(async () => {
    try {
      const response = await fetch('/api/feeds/recent');
      if (!response.ok) {
        throw new Error('Falha ao carregar feeds recentes');
      }
      const payload = (await response.json()) as RecentFeedsResponse;
      if (Array.isArray(payload.feeds)) {
        setRecentFeeds(payload.feeds);
      }
    } catch {
      // Fallback silencioso: mantém lista atual.
    }
  }, []);

  useEffect(() => {
    void loadFeeds();
  }, [loadFeeds]);

  const registerRecentFeed = useCallback((recent: RecentFeed) => {
    setRecentFeeds((current) => {
      const withoutDuplicated = current.filter((item) => item.sourceUrl !== recent.sourceUrl);
      return [recent, ...withoutDuplicated].slice(0, MAX_ITEMS);
    });
  }, []);

  const clearRecentFeeds = useCallback(async () => {
    try {
      const response = await fetch('/api/feeds/recent', { method: 'DELETE' });
      if (!response.ok && response.status !== 204) {
        throw new Error('Falha ao limpar feeds recentes');
      }
      setRecentFeeds([]);
    } catch {
      // Ignora erros para manter UI responsiva.
    }
  }, []);

  return {
    recentFeeds,
    registerRecentFeed,
    clearRecentFeeds,
    reload: loadFeeds,
  };
};
