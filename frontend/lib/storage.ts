const CLIENT_ID_STORAGE_KEY = "ticketeer-client-id";
const QUEUE_TOKEN_STORAGE_KEY_PREFIX = "ticketeer-queue-token:";

export function getOrCreateClientId(): string {
  if (typeof window === "undefined") return "";

  const existing = window.localStorage.getItem(CLIENT_ID_STORAGE_KEY);
  if (existing) return existing;

  const created = crypto.randomUUID();
  window.localStorage.setItem(CLIENT_ID_STORAGE_KEY, created);
  return created;
}

export function saveQueueToken(eventId: number, queueToken: string): void {
  if (typeof window === "undefined") return;
  window.sessionStorage.setItem(`${QUEUE_TOKEN_STORAGE_KEY_PREFIX}${eventId}`, queueToken);
}

export function getQueueToken(eventId: number): string {
  if (typeof window === "undefined") return "";
  return window.sessionStorage.getItem(`${QUEUE_TOKEN_STORAGE_KEY_PREFIX}${eventId}`) ?? "";
}

export function clearQueueToken(eventId: number): void {
  if (typeof window === "undefined") return;
  window.sessionStorage.removeItem(`${QUEUE_TOKEN_STORAGE_KEY_PREFIX}${eventId}`);
}
