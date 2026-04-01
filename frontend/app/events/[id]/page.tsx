import Link from "next/link";
import { notFound } from "next/navigation";
import { StatusBadge } from "@/components/status-badge";
import { fetchEvent } from "@/lib/api";
import { formatDateTime } from "@/lib/format";

type EventDetailPageProps = {
  params: {
    id: string;
  };
};

export const dynamic = "force-dynamic";

export default async function EventDetailPage({ params }: EventDetailPageProps) {
  const event = await fetchEvent(Number(params.id));

  if (!event) {
    notFound();
  }

  const canEnterQueue = event.status === "OPEN";

  return (
    <section className="space-y-6">
      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <div className="flex items-start justify-between gap-4">
          <div className="space-y-3">
            <h1 className="text-3xl font-bold text-slate-900">{event.title}</h1>
            <p className="text-sm text-slate-600">{event.description}</p>
          </div>
          <StatusBadge value={event.status} />
        </div>

        <div className="mt-6 grid gap-4 text-sm text-slate-700 md:grid-cols-2">
          <div>
            <p className="font-semibold text-slate-900">공연장</p>
            <p>{event.venue}</p>
          </div>
          <div>
            <p className="font-semibold text-slate-900">공연 일시</p>
            <p>{formatDateTime(event.event_at)}</p>
          </div>
          <div>
            <p className="font-semibold text-slate-900">예매 시작</p>
            <p>{formatDateTime(event.booking_open_at)}</p>
          </div>
          <div>
            <p className="font-semibold text-slate-900">예매 종료</p>
            <p>{formatDateTime(event.booking_close_at)}</p>
          </div>
        </div>
      </div>

      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-slate-900">예매 안내</h2>
        <p className="mt-2 text-sm text-slate-600">
          예매 오픈 후 대기열에 진입할 수 있습니다. MVP 단계에서는 대기열 진입 후 좌석 선택,
          좌석 홀드, 예매 확정까지 이어집니다.
        </p>

        <div className="mt-6 flex gap-3">
          <Link
            href={canEnterQueue ? `/queue/${event.id}` : "#"}
            className={[
              "inline-flex items-center rounded-xl px-4 py-2 text-sm font-semibold",
              canEnterQueue
                ? "bg-brand-600 text-white hover:bg-brand-700"
                : "cursor-not-allowed bg-slate-200 text-slate-500",
            ].join(" ")}
            aria-disabled={!canEnterQueue}
          >
            대기열 입장
          </Link>

          <Link
            href="/events"
            className="inline-flex items-center rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50"
          >
            목록으로
          </Link>
        </div>
      </div>
    </section>
  );
}
