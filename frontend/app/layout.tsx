import "./globals.css";
import type { Metadata } from "next";
import { AppHeader } from "@/components/app-header";
import { QueryProvider } from "@/components/providers/query-provider";

export const metadata: Metadata = {
  title: "Ticketeer",
  description: "Ticket booking MVP frontend",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ko">
      <body>
        <QueryProvider>
          <div className="min-h-screen bg-slate-50">
            <AppHeader />
            <main className="mx-auto max-w-6xl px-6 py-8">{children}</main>
          </div>
        </QueryProvider>
      </body>
    </html>
  );
}
