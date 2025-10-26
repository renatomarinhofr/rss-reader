type StatusTone = 'info' | 'success' | 'error';

type StatusBannerProps = {
  message: string;
  tone?: StatusTone;
  onDismiss?: () => void;
};

export const StatusBanner = ({ message, tone = 'info', onDismiss }: StatusBannerProps) => {
  if (!message) {
    return null;
  }

  return (
    <div className={`status-banner status-banner--${tone}`} role={tone === 'error' ? 'alert' : 'status'}>
      <span>{message}</span>
      {onDismiss && (
        <button type="button" className="status-banner__dismiss" onClick={onDismiss} aria-label="Fechar aviso">
          x
        </button>
      )}
    </div>
  );
};
