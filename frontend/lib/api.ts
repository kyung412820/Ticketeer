import type { Event } from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    cache: "no-store",
  });

  if (!response.ok) {
    return Promise.reject(new Error(`API request failed: ${response.status}`));
  }

  const json = (await response.json()) as { data: T };
  return json.data;
}

export async function fetchEvents(): Promise<Event[]> {
  return request<Event[]>("/events");
}

export async function fetchEvent(id: number): Promise<Event | null> {
  try {
    return await request<Event>(`/events/${id}`);
  } catch {
    return null;
  }
}
