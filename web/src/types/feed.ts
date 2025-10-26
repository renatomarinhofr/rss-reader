export type FeedItem = {
  title: string;
  link: string;
  description: string;
  publishedAt: string;
};

export type FeedResponse = {
  sourceUrl: string;
  title: string;
  description: string;
  link: string;
  fetchedAt: string;
  items: FeedItem[];
};

export type RecentFeed = {
  sourceUrl: string;
  title: string;
  description: string;
  link: string;
  fetchedAt: string;
};
