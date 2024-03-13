import { CrawlerState } from "@/hooks/useCrawlerState";
import { createContext } from "react";

const CrawlerContext = createContext<CrawlerState>({
    isRunning: false,
    inQueue: 0,
    processed: 0,
    lastRunDuration: 0,
    lastRunStartedAt: 0,
    lastRunEndedAt: 0,
});

export default CrawlerContext;