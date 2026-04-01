"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { enterQueue } from "@/lib/api-client";
import type { QueueStatusResponse } from "@/lib/types";
import { getOrCreateClientId, getQueueToken, saveQueueToken, clearQueueToken } from "@/lib/storage";
import { ErrorBanner } from "@/components/error-banner";
import { LoadingState } from "@/components/loading-state";
import { useQueueStatusQuery } from "@/hooks/use-queue-status-query";

type QueuePageProps = {
  params: {
    eventId: string;
  };
};

export default function QueuePage({ params }: QueuePageProps) {
  const eventId = Number(params.eventId);

  const [clientId, setClientId] = useState<string>("");
  const [queueToken, setQueueToken] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const [isEnteringQueue, setIsEnteringQueue] = useState<boolean>(false);

  const hasEnteredRef = useRef(false);

  useEffect(() => {
    const cid = getOrCreateClientId();
    setClientId(cid);

    const existingQueueToken = getQueueToken(eventId);
    if (existingQueueToken) {
      setQueueToken(existingQueueToken);
      hasEnteredRef.current = true;
      setLoading(false);
      return;
    }

    const run = async () => {
      if (hasEnteredRef.current) return;

      hasEnteredRef.current = true;
      setIsEnteringQueue(true);
      setError("");
      setLoading(true);

      try {
        const entered = await enterQueue({ event_id: eventId, client_id: cid });
        setQueueToken(entered.queue_token);
        saveQueueToken(eventId, entered.queue_token);
      } catch (err) {
        const message = err instanceof Error ? err.message : "대기열 진입 중 오류가 발생했습니다.";
        setError(message);
        setLoading(false);
      } finally {
        setIsEnteringQueue(false);
      }
    };

    void run();
  }, [eventId]);

  const queueStatusQuery = useQueueStatusQuery(queueToken, Boolean(queueToken));
  const queueState = queueStatusQuery.data as QueueStatusResponse | undefined;

  useEffect(() => {
    if (!queueToken) return;

    if (queueStatusQuery.isSuccess) {
      setLoading(false);

      if (queueState?.status === "EXPIRED") {
        clearQueueToken(eventId);
      }
    }

    if (queueStatusQuery.isError) {
      setLoading(false);
      setError(queueStatusQuery.error instanceof Error ? queueStatusQuery.error.message : "대기열 상태 조회 중 오류가 발생했습니다.");
      clearQueueToken(eventId);
    }
  }, [eventId, queueState?.status, queueStatusQuery.error, queueStatusQuery.isError, queueStatusQuery.isSuccess, queueToken]);

  const bookingHref = useMemo(() => {
    if (!queueToken) return "#";
    return `/booking/${eventId}?queueToken=${encodeURIComponent(queueToken)}&clientId=${encodeURIComponent(clientId)}`;
  }, [clientId, eventId, queueToken]);

  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-slate-900">대기열</h1>
        <p className="mt-2 text-sm text-slate-600">
          대기 상태를 주기적으로 갱신합니다.
        </p>
      </div>

      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        {loading || queueStatusQuery.isLoading ? (
          <LoadingState
            message={isEnteringQueue ? "대기열에 진입하는 중입니다..." : "대기열 상태를 확인하는 중입니다..."}
          />
        ) : null}

        {!loading && error ? <ErrorBanner message={error} /> : null}

        {!loading && !error && queueState ? (
          <div className="space-y-4">
            <div className="grid gap-4 md:grid-cols-3" aria-live="polite">
              <Info label="Queue Token" value={queueState.queue_token} />
              <Info label="상태" value={queueState.status} />
              <Info label="순번" value={String(queueState.position ?? 0)} />
            </div>

            <div aria-live="polite" className="rounded-xl bg-slate-50 p-4 text-sm text-slate-700">
              {queueState.status === "WAITING" && "대기 중입니다. 자동으로 상태를 갱신합니다."}
              {queueState.status === "READY" && "입장 가능합니다. 좌석 선택 페이지로 이동하세요."}
              {queueState.status === "EXPIRED" && "입장 시간이 만료되었습니다. 다시 대기열에 진입해주세요."}
            </div>

            <div className="flex gap-3">
              <Link
                href={queueState.status === "READY" ? bookingHref : "#"}
                className={[
                  "inline-flex items-center rounded-xl px-4 py-2 text-sm font-semibold",
                  queueState.status === "READY"
                    ? "bg-brand-600 text-white hover:bg-brand-700"
                    : "cursor-not-allowed bg-slate-200 text-slate-500",
                ].join(" ")}
                aria-disabled={queueState.status !== "READY"}
              >
                좌석 선택으로 이동
              </Link>

              <Link
                href={`/events/${eventId}`}
                className="inline-flex items-center rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50"
              >
                상세로 돌아가기
              </Link>
            </div>
          </div>
        ) : null}
      </div>
    </section>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-4">
      <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">{label}</p>
      <p className="mt-2 break-all text-sm font-medium text-slate-900">{value}</p>
    </div>
  );
}
