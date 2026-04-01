import Link from "next/link";
import { memo, useMemo } from "react";
import type { Seat } from "@/lib/types";

type SeatSelectionSummaryProps = {
  eventId: number;
  selectedSeat: Seat | null;
  holdExpiresAt: string;
  isHolding: boolean;
  isBooking: boolean;
  hasBookingResult: boolean;
  onHold: () => void;
  onBooking: () => void;
};

export const SeatSelectionSummary = memo(function SeatSelectionSummary({
  eventId,
  selectedSeat,
  holdExpiresAt,
  isHolding,
  isBooking,
  hasBookingResult,
  onHold,
  onBooking,
}: SeatSelectionSummaryProps) {
  const remainingText = useMemo(() => {
    return holdExpiresAt ? new Date(holdExpiresAt).toLocaleString("ko-KR") : "-";
  }, [holdExpiresAt]);

  const holdDisabled = !selectedSeat || selectedSeat.status !== "AVAILABLE" || isHolding || isBooking || hasBookingResult;
  const bookingDisabled = !selectedSeat || selectedSeat.status !== "AVAILABLE" || isBooking || isHolding || hasBookingResult;

  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <h2 className="text-lg font-semibold text-slate-900">선택 좌석</h2>

      <div aria-live="polite">
        {selectedSeat ? (
          <div className="mt-4 space-y-2 text-sm text-slate-700">
            <p><span className="font-semibold text-slate-900">좌석:</span> {selectedSeat.seat_no}</p>
            <p><span className="font-semibold text-slate-900">구역:</span> {selectedSeat.section}</p>
            <p><span className="font-semibold text-slate-900">가격:</span> {selectedSeat.price.toLocaleString("ko-KR")}원</p>
            <p><span className="font-semibold text-slate-900">현재 상태:</span> {selectedSeat.status}</p>
            <p><span className="font-semibold text-slate-900">홀드 만료:</span> {remainingText}</p>
          </div>
        ) : (
          <p className="mt-4 text-sm text-slate-600">좌석을 선택하세요.</p>
        )}
      </div>

      <div className="mt-6 flex flex-col gap-3">
        <button
          type="button"
          onClick={onHold}
          disabled={holdDisabled}
          aria-disabled={holdDisabled}
          className="rounded-xl bg-brand-600 px-4 py-2 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:bg-slate-300"
        >
          {isHolding ? "좌석 홀드 중..." : "좌석 홀드"}
        </button>
        <button
          type="button"
          onClick={onBooking}
          disabled={bookingDisabled}
          aria-disabled={bookingDisabled}
          className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold text-slate-700 disabled:cursor-not-allowed disabled:bg-slate-100"
        >
          {isBooking ? "예매 확정 중..." : hasBookingResult ? "예매 완료" : "예매 확정"}
        </button>
        <Link
          href={`/queue/${eventId}`}
          className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-center text-sm font-semibold text-slate-700 hover:bg-slate-50"
        >
          대기열로 돌아가기
        </Link>
      </div>
    </div>
  );
});
