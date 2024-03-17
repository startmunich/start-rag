"use client";
import { Package, ScanSearch } from "lucide-react";
import { TooltipProvider } from "@radix-ui/react-tooltip";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "./ui/resizable";
import { CrawlerSwitcher } from "./crawler-switcher";
import { Separator } from "@radix-ui/react-select";
import { Nav } from "./nav";
import { cn } from "@/lib/utils";
import { useContext } from "react";
import AppStateContext from "@/context/appStateContext";
import { usePathname } from "next/navigation";

interface DashboardProps {
  selectableCrawler: {
    label: string;
    name: string;
    icon: React.ReactNode;
  }[];
  defaultLayout: number[];
  children: React.ReactNode;
}

export default function Dashboard({
  selectableCrawler,
  defaultLayout,
  children,
}: DashboardProps) {
  const { crawler, pageCount } = useContext(AppStateContext);
  const pathname = usePathname();

  return (
    <TooltipProvider delayDuration={0}>
      <ResizablePanelGroup
        direction="horizontal"
        onLayout={(sizes: number[]) => {
          document.cookie = `react-resizable-panels:layout=${JSON.stringify(
            sizes,
          )}`;
        }}
        className="h-[100vh] items-stretch"
      >
        <ResizablePanel
          defaultSize={defaultLayout[0]}
          collapsible={false}
          minSize={15}
          maxSize={20}
        >
          <div
            className={cn(
              "flex h-[52px] items-center justify-center px-2",
            )}
          >
            <CrawlerSwitcher isCollapsed={false} crawler={selectableCrawler} />
          </div>
          <Separator />
          <Nav
            isCollapsed={false}
            links={[
              {
                title: "Crawler",
                href: "/crawler",
                label: crawler.isRunning ? "Running" : "Stopped",
                icon: ScanSearch,
                isSelected: pathname === "/crawler",
              },
              {
                title: "Workspace",
                href: "/workspace",
                label: pageCount.count.toString(),
                icon: Package,
                isSelected: pathname === "/workspace",
              },
            ]}
          />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={defaultLayout[1]} minSize={30}>
            {children}
        </ResizablePanel>
      </ResizablePanelGroup>
    </TooltipProvider>
  );
}
