import Link from "next/link";

export function AppHeader() {
  return (
    <header className="border-b border-slate-200 bg-white">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/events" className="text-xl font-bold text-brand-700">
          TICKETEER
        </Link>
        <nav className="flex items-center gap-4 text-sm text-slate-600">
          <Link href="/events" className="hover:text-slate-900">
            공연 목록
          </Link>
        </nav>
      </div>
    </header>
  );
}
