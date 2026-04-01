export const queryKeys = {
  seats: (eventId: number) => ["seats", eventId] as const,
  queueStatus: (queueToken: string) => ["queue-status", queueToken] as const,
};
