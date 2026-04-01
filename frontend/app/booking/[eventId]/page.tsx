"use client";

import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { createBooking, holdSeat } from "@/lib/api-client";
import type { BookingResponse, Seat } from "@/lib/types";
import { SeatBadge } from "@/components/status-badge";
import { clearQueueToken, getOrCreateClientId, getQueueToken } from "@/lib/storage";
import { SeatSectionGrid } from "@/components/seat-section-grid";
import { LoadingState } from "@/components/loading-state";
import { ErrorBanner } from "@/components/error-banner";
import { BookingResultCard } from "@/components/booking-result-card";
import { SeatSelectionSummary } from "@/components/seat-selection-summary";
import { useSeatsQuery } from "@/hooks/use-seats-query";
import { queryKeys } from "@/lib/query-keys";

type BookingPageProps = {
  params: {
    eventId: string;
  };
};

export default function BookingPage({ params }: BookingPageProps) {
  const eventId = Number(params.eventId);
  const searchParams = useSearchParams();
  const queueTokenFromQuery = searchParams.get("queueToken") ?? "";
  const clientIdFromQuery = searchParams.get("clientId") ?? "";

  const [resolvedQueueToken, setResolvedQueueToken] = useState<string>("");
  const [resolvedClientId, setResolvedClientId] = useState<string>("");
  const [selectedSeat, setSelectedSeat] = useState<Seat | null>(null);
  const [heldSeatSnapshot, setHeldSeatSnapshot] = useState<Seat | null>(null);
  const [holdExpiresAt, setHoldExpiresAt] = useState<string>("");
  const [bookingResult, setBookingResult] = useState<BookingResponse | null>(null);
  const [message, setMessage] = useState<string>("");
  const [isHolding, setIsHolding] = useState<boolean>(false);
  const [isBooking, setIsBooking] = useState<boolean>(false);

  const queryClient = useQueryClient();

  useEffect(() => {
    const clientId = clientIdFromQuery || getOrCreateClientId();
    const queueToken = queueTokenFromQuery || getQueueToken(eventId);

    setResolvedClientId(clientId);
    setResolvedQueueToken(queueToken);

    if (!queueToken || !clientId) {
      setMessage("유효한 접근이 아닙니다. 대기열부터 다시 진행해주세요.");
    }
  }, [clientIdFromQuery, eventId, queueTokenFromQuery]);

  const seatsQuery = useSeatsQuery(eventId, Boolean(resolvedQueueToken && resolvedClientId));
  const seats = seatsQuery.data?.seats ?? [];

  useEffect(() => {
    setSelectedSeat((current) => {
      if (!current) return current;
      const latest = seats.find((seat) => seat.id === current.id);
      return latest ?? null;
    });
  }, [seats]);

  const remainingText = useMemo(() => {
    if (!holdExpiresAt) return "-";
    return new Date(holdExpiresAt).toLocaleString("ko-KR");
  }, [holdExpiresAt]);

  const resultSeatLabel = useMemo(() => {
    if (bookingResult && heldSeatSnapshot) return heldSeatSnapshot.seat_no;
    return heldSeatSnapshot?.seat_no ?? selectedSeat?.seat_no ?? "-";
  }, [bookingResult, heldSeatSnapshot, selectedSeat]);

  const resultSeatSection = useMemo(() => {
    if (bookingResult && heldSeatSnapshot) return heldSeatSnapshot.section;
    return heldSeatSnapshot?.section ?? selectedSeat?.section ?? "-";
  }, [bookingResult, heldSeatSnapshot, selectedSeat]);

  const refetchSeats = async () => {
    await queryClient.invalidateQueries({ queryKey: queryKeys.seats(eventId) });
  };

  const handleHold = async () => {
    if (!selectedSeat || isHolding || isBooking) return;

    try {
      setIsHolding(true);
      setMessage("");
      const result = await holdSeat(selectedSeat.id, {
        event_id: eventId,
        client_id: resolvedClientId,
        queue_token: resolvedQueueToken,
      });
      setHoldExpiresAt(result.hold_expires_at);
      setHeldSeatSnapshot(selectedSeat);
      setMessage("좌석 홀드에 성공했습니다.");
      await refetchSeats();
    } catch (err) {
      const msg = err instanceof Error ? err.message : "좌석 홀드 중 오류가 발생했습니다.";
      setMessage(msg);
      await refetchSeats();
    } finally {
      setIsHolding(false);
    }
  };

  const handleBooking = async () => {
    if (!selectedSeat || isBooking || isHolding || bookingResult) return;

    try {
      setIsBooking(true);
      setMessage("");
      const result = await createBooking({
        event_id: eventId,
        seat_id: selectedSeat.id,
        client_id: resolvedClientId,
        queue_token: resolvedQueueToken,
      });
      setBookingResult(result);
      setHeldSeatSnapshot(selectedSeat);
      setMessage("예매가 확정되었습니다.");
      clearQueueToken(eventId);
      await refetchSeats();
    } catch (err) {
      const msg = err instanceof Error ? err.message : "예매 확정 중 오류가 발생했습니다.";
      setMessage(msg);
      await refetchSeats();
    } finally {
      setIsBooking(false);
    }
  };

  const selectionDisabled = isHolding || isBooking || Boolean(bookingResult);

  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-slate-900">좌석 선택</h1>
        <p className="mt-2 text-sm text-slate-600">
          구역별 좌석을 선택하고 홀드한 뒤 예매를 확정합니다.
        </p>
      </div>

      {message && !bookingResult ? <ErrorBanner message={message} /> : null}

      <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="mb-4 flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold text-slate-900">좌석 목록</h2>
              <p className="mt-1 text-sm text-slate-500">구역별로 좌석을 묶어서 표시합니다.</p>
            </div>
            <div className="flex gap-2">
              <SeatBadge value="AVAILABLE" />
              <SeatBadge value="HELD" />
              <SeatBadge value="BOOKED" />
            </div>
          </div>

          {seatsQuery.isLoading ? (
            <LoadingState message="좌석 정보를 불러오는 중입니다..." />
          ) : seatsQuery.isError ? (
            <ErrorBanner message={seatsQuery.error instanceof Error ? seatsQuery.error.message : "좌석 조회 중 오류가 발생했습니다."} />
          ) : (
            <SeatSectionGrid
              seats={seats}
              selectedSeat={selectedSeat}
              disabled={selectionDisabled}
              onSelect={setSelectedSeat}
            />
          )}
        </div>

        <aside className="space-y-4">
          <SeatSelectionSummary
            eventId={eventId}
            selectedSeat={selectedSeat}
            holdExpiresAt={holdExpiresAt}
            isHolding={isHolding}
            isBooking={isBooking}
            hasBookingResult={Boolean(bookingResult)}
            onHold={handleHold}
            onBooking={handleBooking}
          />

          {bookingResult ? (
            <BookingResultCard
              bookingResult={bookingResult}
              eventId={eventId}
              seatLabel={resultSeatLabel}
              seatSection={resultSeatSection}
              message={message}
            />
          ) : null}
        </aside>
      </div>
    </section>
  );
}
