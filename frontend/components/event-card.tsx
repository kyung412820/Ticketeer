import Link from "next/link";
import { StatusBadge } from "@/components/status-badge";
import type { Event } from "@/lib/types";

type EventCardProps = {
  event: Event;
  eventAtText: string;
  bookingOpenText: string;
};

export function EventCard({ event, eventAtText, bookingOpenText }: EventCardProps) {
  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h2 className="text-xl font-semibold text-slate-900">{event.title}</h2>
          <p className="mt-2 text-sm text-slate-600">{event.venue}</p>
        </div>
        <StatusBadge value={event.status} />
      </div>

      <div className="mt-4 space-y-2 text-sm text-slate-700">
        <p><span className="font-semibold text-slate-900">공연일시:</span> {eventAtText}</p>
        <p><span className="font-semibold text-slate-900">예매시작:</span> {bookingOpenText}</p>
      </div>

      <Link
        href={`/events/${event.id}`}
        className="mt-6 inline-flex rounded-xl bg-brand-600 px-4 py-2 text-sm font-semibold text-white hover:bg-brand-700"
      >
        상세 보기
      </Link>
    </article>
  );
}
