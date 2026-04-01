"use client";

import { memo, useEffect, useMemo, useRef } from "react";
import type { Seat } from "@/lib/types";
import { SeatBadge } from "@/components/status-badge";

type SeatSectionGridProps = {
  seats: Seat[];
  selectedSeat: Seat | null;
  disabled: boolean;
  onSelect: (seat: Seat) => void;
};

type SeatSectionProps = {
  section: string;
  seats: Seat[];
  selectedSeatId?: number;
  disabled: boolean;
  onSelect: (seat: Seat) => void;
};

export const SeatSectionGrid = memo(function SeatSectionGrid({
  seats,
  selectedSeat,
  disabled,
  onSelect,
}: SeatSectionGridProps) {
  const grouped = useMemo(() => groupBySection(seats), [seats]);

  return (
    <div className="space-y-6" role="list" aria-label="구역별 좌석 목록">
      {grouped.map(([section, sectionSeats]) => (
        <SeatSection
          key={section}
          section={section}
          seats={sectionSeats}
          selectedSeatId={selectedSeat?.id}
          disabled={disabled}
          onSelect={onSelect}
        />
      ))}
    </div>
  );
});

const SeatSection = memo(function SeatSection({
  section,
  seats,
  selectedSeatId,
  disabled,
  onSelect,
}: SeatSectionProps) {
  return (
    <section className="space-y-3" aria-label={`${section} 구역 좌석 목록`}>
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-base font-semibold text-slate-900">{section} 구역</h3>
          <p className="text-xs text-slate-500">{seats.length}개 좌석</p>
        </div>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3" role="list">
        {seats.map((seat) => (
          <SeatCard
            key={seat.id}
            seat={seat}
            isSelected={selectedSeatId === seat.id}
            disabled={disabled}
            onSelect={onSelect}
          />
        ))}
      </div>
    </section>
  );
});

type SeatCardProps = {
  seat: Seat;
  isSelected: boolean;
  disabled: boolean;
  onSelect: (seat: Seat) => void;
};

const SeatCard = memo(function SeatCard({
  seat,
  isSelected,
  disabled,
  onSelect,
}: SeatCardProps) {
  const isAvailable = seat.status === "AVAILABLE";
  const isDisabled = !isAvailable || disabled;
  const buttonRef = useRef<HTMLButtonElement | null>(null);

  useEffect(() => {
    if (isSelected && buttonRef.current) {
      buttonRef.current.focus();
    }
  }, [isSelected]);

  return (
    <button
      ref={buttonRef}
      type="button"
      onClick={() => onSelect(seat)}
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-pressed={isSelected}
      aria-current={isSelected ? "true" : undefined}
      aria-label={`${seat.section} 구역 ${seat.seat_no} 좌석, ${seat.price.toLocaleString("ko-KR")}원, 현재 상태 ${seat.status}`}
      className={[
        "rounded-xl border p-4 text-left transition focus:outline-none focus:ring-2 focus:ring-brand-200 focus:ring-offset-2",
        isSelected ? "border-brand-600 ring-2 ring-brand-100" : "border-slate-200",
        !isDisabled ? "bg-white hover:bg-slate-50" : "cursor-not-allowed bg-slate-50 opacity-80",
      ].join(" ")}
    >
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-slate-900">{seat.seat_no}</p>
          <p className="mt-1 text-xs text-slate-600">{seat.section} 구역</p>
        </div>
        <SeatBadge value={seat.status} />
      </div>
      <p className="mt-3 text-sm font-medium text-slate-800">
        {seat.price.toLocaleString("ko-KR")}원
      </p>
    </button>
  );
});

function groupBySection(seats: Seat[]): Array<[string, Seat[]]> {
  const map = new Map<string, Seat[]>();

  for (const seat of seats) {
    const current = map.get(seat.section) ?? [];
    current.push(seat);
    map.set(seat.section, current);
  }

  return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0]));
}
