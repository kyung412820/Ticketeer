type StatusBadgeProps = {
  value: string;
};

export function StatusBadge({ value }: StatusBadgeProps) {
  const className = getBadgeClassName(value);

  return (
    <span className={["inline-flex rounded-full px-3 py-1 text-xs font-semibold", className].join(" ")}>
      {value}
    </span>
  );
}

export function SeatBadge({ value }: StatusBadgeProps) {
  const className = getBadgeClassName(value);

  return (
    <span className={["inline-flex rounded-full px-2.5 py-1 text-[11px] font-semibold", className].join(" ")}>
      {value}
    </span>
  );
}

function getBadgeClassName(value: string): string {
  switch (value) {
    case "OPEN":
    case "AVAILABLE":
    case "CONFIRMED":
    case "READY":
      return "bg-green-100 text-green-700";
    case "OPEN_PENDING":
    case "WAITING":
      return "bg-blue-100 text-blue-700";
    case "HELD":
      return "bg-orange-100 text-orange-700";
    case "BOOKED":
    case "CLOSED":
    case "EXPIRED":
      return "bg-slate-200 text-slate-700";
    default:
      return "bg-slate-100 text-slate-700";
  }
}
