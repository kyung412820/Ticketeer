type LoadingStateProps = {
  message: string;
};

export function LoadingState({ message }: LoadingStateProps) {
  return (
    <p data-testid="loading-state" aria-live="polite" className="text-sm text-slate-600">
      {message}
    </p>
  );
}
