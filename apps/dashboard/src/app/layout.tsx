"use client";
import { Inter } from "next/font/google";
import "./globals.css";
import Dashboard from "@/components/dashboard";
import { NotionLogoIcon } from "@radix-ui/react-icons";
import AppStateContext from "@/context/appStateContext";
import useAppState from "@/hooks/useAppState";

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
  const appState = useAppState(process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!);
  return (
    <html lang="en">
      <body className={inter.className}>
        <AppStateContext.Provider value={appState}>
          <Dashboard
            selectableCrawler={crawler}
            defaultLayout={[265, 440, 655]}
          >
            {children}
          </Dashboard>
        </AppStateContext.Provider>
      </body>
    </html>
  );
}
