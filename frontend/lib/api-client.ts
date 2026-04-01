import type {
  BookingResponse,
  HoldSeatResponse,
  QueueEnterResponse,
  QueueStatusResponse,
  SeatsResponse,
} from "@/lib/types";
import { getUserFriendlyErrorMessage } from "@/lib/error-message";

type RequestOptions = RequestInit & {
  timeoutMs?: number;
  retryCount?: number;
  retryDelayMs?: number;
};

const API_BASE_URL =
  typeof window === "undefined"
    ? process.env.INTERNAL_API_BASE_URL ??
      process.env.NEXT_PUBLIC_API_BASE_URL ??
      "http://backend:8080/api"
    : process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function request<T>(path: string, options?: RequestOptions): Promise<T> {
  const {
    timeoutMs = 5000,
    retryCount = 0,
    retryDelayMs = 300,
    ...init
  } = options ?? {};

  let lastError: Error | null = null;

  for (let attempt = 0; attempt <= retryCount; attempt += 1) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

    try {
      const response = await fetch(`${API_BASE_URL}${path}`, {
        ...init,
        signal: controller.signal,
        headers: {
          "Content-Type": "application/json",
          ...(init?.headers ?? {}),
        },
      });

      const json = await response.json().catch(() => null);
      const errorCode = json?.error?.code as string | undefined;
      const fallbackMessage = json?.error?.message as string | undefined;

      if (!response.ok) {
        throw new Error(getUserFriendlyErrorMessage(errorCode, fallbackMessage));
      }

      return json.data as T;
    } catch (error) {
      if (error instanceof DOMException && error.name === "AbortError") {
        lastError = new Error("요청 시간이 초과되었습니다. 잠시 후 다시 시도해주세요.");
      } else if (error instanceof Error) {
        lastError = error.message.includes("Failed to fetch")
          ? new Error("서버에 연결할 수 없습니다. 백엔드 실행 상태를 확인해주세요.")
          : error;
      } else {
        lastError = new Error("요청 처리 중 오류가 발생했습니다.");
      }

      if (attempt < retryCount) {
        await sleep(retryDelayMs);
        continue;
      }
    } finally {
      clearTimeout(timeoutId);
    }
  }

  throw lastError ?? new Error("요청 처리 중 오류가 발생했습니다.");
}

export async function enterQueue(body: { event_id: number; client_id: string }) {
  return request<QueueEnterResponse>("/queue/enter", {
    method: "POST",
    body: JSON.stringify(body),
    timeoutMs: 5000,
  });
}

export async function getQueueStatus(queueToken: string) {
  return request<QueueStatusResponse>(`/queue/status/${queueToken}`, {
    timeoutMs: 3000,
    retryCount: 1,
    retryDelayMs: 250,
  });
}

export async function fetchSeats(eventId: number) {
  return request<SeatsResponse>(`/events/${eventId}/seats`, {
    timeoutMs: 4000,
    retryCount: 1,
    retryDelayMs: 250,
  });
}

export async function holdSeat(
  seatId: number,
  body: { event_id: number; client_id: string; queue_token: string },
) {
  return request<HoldSeatResponse>(`/seats/${seatId}/hold`, {
    method: "POST",
    body: JSON.stringify(body),
    timeoutMs: 5000,
  });
}

export async function createBooking(body: {
  event_id: number;
  seat_id: number;
  client_id: string;
  queue_token: string;
}) {
  return request<BookingResponse>("/bookings", {
    method: "POST",
    body: JSON.stringify(body),
    timeoutMs: 5000,
  });
}
