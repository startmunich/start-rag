"use client";
import { Inter } from "next/font/google";
import "./globals.css";
import Dashboard from "@/components/dashboard";
import { ChatBubbleIcon } from "@radix-ui/react-icons";
import AppStateContext from "@/context/appStateContext";
import useAppState from "@/hooks/useAppState";

const inter = Inter({ subsets: ["latin"] });

export const crawler = [
  {
    label: "START RAG",
    name: "START RAG",
    icon: <ChatBubbleIcon />,
  },
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
