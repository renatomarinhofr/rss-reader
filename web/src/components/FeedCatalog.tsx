import type { FeedCategory, FeedSource } from '../data/sources';

type FeedCatalogProps = {
  categories: FeedCategory[];
  activeCategoryId: string;
  selectedFeedUrl?: string | null;
  onCategorySelect: (categoryId: string) => void;
  onFeedSelect: (categoryId: string, feed: FeedSource) => void;
};

export const FeedCatalog = ({
  categories,
  activeCategoryId,
  selectedFeedUrl,
  onCategorySelect,
  onFeedSelect,
}: FeedCatalogProps) => {
  if (categories.length === 0) {
    return null;
  }

  const activeCategory =
    categories.find((category) => category.id === activeCategoryId) ?? categories[0];

  return (
    <section className="feed-catalog">
      <header className="feed-catalog__header">
        <div>
          <h2>Descubra fontes brasileiras</h2>
          <p>
            Selecione uma categoria para explorar fontes confiáveis. O conteúdo é carregado
            automaticamente e você pode salvar os links favoritos.
          </p>
        </div>
      </header>

      <nav className="feed-catalog__categories" aria-label="Categorias de feeds">
        {categories.map((category) => {
          const isActive = category.id === activeCategory.id;
          return (
            <button
              key={category.id}
              type="button"
              className={`feed-catalog__category ${isActive ? 'is-active' : ''}`}
              onClick={() => onCategorySelect(category.id)}
            >
              <span>{category.name}</span>
            </button>
          );
        })}
      </nav>

      <p className="feed-catalog__description">{activeCategory.description}</p>

      <div className="feed-catalog__grid">
        {activeCategory.feeds.map((feed) => {
          const isSelected = selectedFeedUrl === feed.url;
          return (
            <article
              key={feed.url}
              className={`feed-catalog__card ${isSelected ? 'is-selected' : ''}`}
            >
              <h3>{feed.name}</h3>
              {feed.description && <p>{feed.description}</p>}
              <button
                type="button"
                className="feed-catalog__cta"
                onClick={() => onFeedSelect(activeCategory.id, feed)}
                aria-label={`Carregar notícias de ${feed.name}`}
              >
                Ler feed
              </button>
            </article>
          );
        })}
      </div>
    </section>
  );
};
