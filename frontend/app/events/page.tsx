import { EventCard } from "@/components/event-card";
import { fetchEvents } from "@/lib/api";
import { formatDateTime } from "@/lib/format";

export const dynamic = "force-dynamic";

export default async function EventsPage() {
  const events = await fetchEvents();

  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-slate-900">공연 목록</h1>
        <p className="mt-2 text-sm text-slate-600">
          예매 상태를 확인하고 상세 페이지로 이동하세요.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {events.map((event) => (
          <EventCard
            key={event.id}
            event={event}
            bookingOpenText={formatDateTime(event.booking_open_at)}
            eventAtText={formatDateTime(event.event_at)}
          />
        ))}
      </div>
    </section>
  );
}
