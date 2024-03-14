"use client";

import * as React from "react";

import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";

interface CrawlerSwitcherProps {
  isCollapsed: boolean;
  crawler: {
    label: string;
    name: string;
    icon: React.ReactNode;
  }[];
}

export function CrawlerSwitcher({
  isCollapsed,
  crawler,
}: CrawlerSwitcherProps) {
  const [selectedAccount, setSelectedAccount] = React.useState<string>(
    crawler[0].name,
  );

  return (
    <Select defaultValue={selectedAccount} onValueChange={setSelectedAccount}>
      <SelectTrigger
        className={cn(
          "flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0",
          isCollapsed &&
            "flex h-9 w-9 shrink-0 items-center justify-center p-0 [&>span]:w-auto [&>svg]:hidden",
        )}
        aria-label="Select account"
      >
        <SelectValue placeholder="Select an account">
          {crawler.find((account) => account.name === selectedAccount)?.icon}
          <span className={cn("ml-2", isCollapsed && "hidden")}>
            {
              crawler.find((account) => account.name === selectedAccount)
                ?.label
            }
          </span>
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {crawler.map((account) => (
          <SelectItem key={account.name} value={account.name}>
            <div className="flex items-center gap-3 [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0 [&_svg]:text-foreground">
              {account.icon}
              {account.name}
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
