import { FormEvent } from 'react';

type FeedFormProps = {
  value: string;
  loading: boolean;
  onChange: (value: string) => void;
  onSubmit: (value: string) => void;
};

export const FeedForm = ({ value, loading, onChange, onSubmit }: FeedFormProps) => {
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    onSubmit(value.trim());
  };

  return (
    <form className="feed-form" onSubmit={handleSubmit}>
      <div className="form-control">
        <label className="form-label" htmlFor="feed-url">
          URL do Feed
        </label>
        <input
          id="feed-url"
          type="url"
          placeholder="https://exemplo.com/rss"
          value={value}
          onChange={(event) => onChange(event.target.value)}
          required
          autoComplete="url"
          spellCheck={false}
        />
      </div>

      <button type="submit" disabled={loading || value.trim().length === 0}>
        {loading ? 'Carregando...' : 'Ler feed'}
      </button>
    </form>
  );
};
