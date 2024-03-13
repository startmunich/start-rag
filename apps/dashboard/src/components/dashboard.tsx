"use client";
import { File, Inbox, Package, ScanSearch } from "lucide-react";
import { TooltipProvider } from "@radix-ui/react-tooltip";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "./ui/resizable";
import { AccountSwitcher } from "./account-switcher";
import { Separator } from "@radix-ui/react-select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import { Nav } from "./nav";
import { cn } from "@/lib/utils";

interface DashboardProps {
  accounts: {
    label: string;
    email: string;
    icon: React.ReactNode;
  }[];
  defaultLayout: number[];
  navCollapsedSize: number;
}

export default function Dashboard({
  accounts,
  defaultLayout,
  navCollapsedSize,
}: DashboardProps) {
  return (
    <TooltipProvider delayDuration={0}>
      <ResizablePanelGroup
        direction="horizontal"
        onLayout={(sizes: number[]) => {
          document.cookie = `react-resizable-panels:layout=${JSON.stringify(
            sizes,
          )}`;
        }}
        className="h-full max-h-[800px] items-stretch"
      >
        <ResizablePanel
          defaultSize={defaultLayout[0]}
          collapsedSize={navCollapsedSize}
          collapsible={false}
          minSize={15}
          maxSize={20}
        >
          <div
            className={cn(
              "flex h-[52px] items-center justify-center px-2",
            )}
          >
            <AccountSwitcher isCollapsed={false} accounts={accounts} />
          </div>
          <Separator />
          <Nav
            isCollapsed={false}
            links={[
              {
                title: "Crawler",
                label: "Running",
                icon: ScanSearch,
                variant: "default",
              },
              {
                title: "Backup",
                label: "20.000",
                icon: Package,
                variant: "ghost",
              },
            ]}
          />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={defaultLayout[1]} minSize={30}>
          <Tabs defaultValue="all">
            <div className="flex items-center px-4 py-2">
              <h1 className="text-xl font-bold">Inbox</h1>
              <TabsList className="ml-auto">
                <TabsTrigger
                  value="all"
                  className="text-zinc-600 dark:text-zinc-200"
                >
                  All
                </TabsTrigger>
                <TabsTrigger
                  value="unread"
                  className="text-zinc-600 dark:text-zinc-200"
                >
                  New
                </TabsTrigger>
              </TabsList>
            </div>
            <Separator />
            <TabsContent value="all" className="m-0">
              Placeholder
            </TabsContent>
            <TabsContent value="unread" className="m-0">
              Placeholder
            </TabsContent>
          </Tabs>
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={defaultLayout[2]}>
          Content of a page
        </ResizablePanel>
      </ResizablePanelGroup>
    </TooltipProvider>
  );
}
