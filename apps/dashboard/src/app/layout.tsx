"use client";
import { Inter } from "next/font/google";
import "./globals.css";
import Dashboard from "@/components/dashboard";
import useCrawlerState from "@/hooks/useCrawlerState";
import CrawlerContext from "@/context/crawlerContext";
import { NotionLogoIcon } from "@radix-ui/react-icons";

const inter = Inter({ subsets: ["latin"] });

export const crawler = [
  {
    label: "START Notion",
    name: "START Notion",
    icon: (
      <NotionLogoIcon className="w-5 h-5" />
    ),
  }
];

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const crawlerState = useCrawlerState(process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!);
  return (
    <html lang="en">
      <body className={inter.className}>
        <CrawlerContext.Provider value={crawlerState}>
          <Dashboard
            crawler={crawler}
            defaultLayout={[265, 440, 655]}
          >
            {children}
          </Dashboard>
        </CrawlerContext.Provider>
      </body>
    </html>
  );
}
