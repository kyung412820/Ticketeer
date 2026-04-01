import React from "react";
type SeatProps = {
  seatId: number;
  seatLabel: string;
  selected: boolean;
  onClick: (seatId: number) => void;
};

export const Seat = ({ seatId, seatLabel, selected, onClick }: SeatProps) => {
  return (
    <div
      role="button"
      aria-pressed={selected}
      aria-label={`좌석 ${seatLabel} ${selected ? "선택됨" : "선택 안 됨"}`}
      tabIndex={0}
      onClick={() => onClick(seatId)}
      onKeyDown={(e) => e.key === 'Enter' && onClick(seatId)}
      className={`seat ${selected ? 'seat-selected' : 'seat-unselected'}`}
    >
      {seatLabel}
    </div>
  );
};
