import Link from "next/link";
import { useEffect, useRef } from "react";
import type { BookingResponse } from "@/lib/types";

type BookingResultCardProps = {
  bookingResult: BookingResponse;
  eventId: number;
  seatLabel: string;
  seatSection: string;
  message?: string;
};

export function BookingResultCard({
  bookingResult,
  eventId,
  seatLabel,
  seatSection,
  message,
}: BookingResultCardProps) {
  const headingRef = useRef<HTMLHeadingElement | null>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, [bookingResult.booking_code]);

  return (
    <div
      className="rounded-2xl border border-green-200 bg-green-50 p-6 shadow-sm"
      role="status"
      aria-live="polite"
      aria-labelledby="booking-complete-heading"
    >
      <h2
        id="booking-complete-heading"
        ref={headingRef}
        tabIndex={-1}
        className="text-lg font-semibold text-slate-900 focus:outline-none"
      >
        예매 완료
      </h2>

      {message && <p className="mt-3 text-sm text-slate-700">{message}</p>}

      <div className="mt-4 space-y-4">
        <div className="rounded-xl bg-white p-4">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">Booking Code</p>
          <p className="mt-2 break-all text-lg font-bold text-slate-900">
            {bookingResult.booking_code}
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          <ResultItem label="공연 ID" value={String(bookingResult.event_id)} />
          <ResultItem label="좌석" value={seatLabel} />
          <ResultItem label="구역" value={seatSection} />
          <ResultItem
            label="예매 시각"
            value={new Date(bookingResult.booked_at).toLocaleString("ko-KR")}
          />
        </div>

        <div className="flex flex-col gap-3">
          <Link
            href="/events"
            className="rounded-xl bg-brand-600 px-4 py-2 text-center text-sm font-semibold text-white hover:bg-brand-700"
          >
            공연 목록으로 이동
          </Link>
          <Link
            href={`/events/${eventId}`}
            className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-center text-sm font-semibold text-slate-700 hover:bg-slate-50"
          >
            공연 상세로 돌아가기
          </Link>
        </div>
      </div>
    </div>
  );
}

function ResultItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-xl bg-white p-4">
      <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">{label}</p>
      <p className="mt-2 text-sm font-medium text-slate-900">{value}</p>
    </div>
  );
}
