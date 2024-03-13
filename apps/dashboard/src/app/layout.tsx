"use client";
import { Inter } from "next/font/google";
import "./globals.css";
import Dashboard from "@/components/dashboard";
import useCrawlerState from "@/hooks/useCrawlerState";
import CrawlerContext from "@/context/crawlerContext";

const inter = Inter({ subsets: ["latin"] });

export const accounts = [
  {
    label: "Alicia Koch",
    email: "alicia@example.com",
    icon: (
      <svg role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <title>Vercel</title>
        <path d="M24 22.525H0l12-21.05 12 21.05z" fill="currentColor" />
      </svg>
    ),
  }
];

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const crawlerState = useCrawlerState("http://localhost:3111");
  return (
    <html lang="en">
      <body className={inter.className}>
        <CrawlerContext.Provider value={crawlerState}>
          <Dashboard
            accounts={accounts}
            defaultLayout={[265, 440, 655]}
            children={children}
          />
        </CrawlerContext.Provider>
      </body>
    </html>
  );
}
