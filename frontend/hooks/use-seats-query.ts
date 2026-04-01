"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchSeats } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";

export function useSeatsQuery(eventId: number, enabled = true) {
  return useQuery({
    queryKey: queryKeys.seats(eventId),
    queryFn: () => fetchSeats(eventId),
    enabled,
    staleTime: 1_000,
  });
}
