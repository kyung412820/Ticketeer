import Link from "next/link";

export default function NotFound() {
  return (
    <section className="rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
      <h1 className="text-2xl font-bold text-slate-900">페이지를 찾을 수 없습니다.</h1>
      <p className="mt-2 text-sm text-slate-600">
        요청한 경로가 없거나, 아직 구현되지 않은 페이지입니다.
      </p>
      <Link
        href="/events"
        className="mt-6 inline-flex rounded-xl bg-brand-600 px-4 py-2 text-sm font-semibold text-white hover:bg-brand-700"
      >
        공연 목록으로
      </Link>
    </section>
  );
}
