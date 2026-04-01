"use client";

import { useQuery } from "@tanstack/react-query";
import { getQueueStatus } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";

export function useQueueStatusQuery(queueToken: string, enabled = true) {
  return useQuery({
    queryKey: queryKeys.queueStatus(queueToken),
    queryFn: () => getQueueStatus(queueToken),
    enabled: enabled && Boolean(queueToken),
    refetchInterval: (query) => {
      const status = query.state.data?.status;
      if (!status) return 2_000;
      if (status === "READY" || status === "EXPIRED") return false;
      return 2_000;
    },
    staleTime: 1_000,
    retry: 1,
  });
}
