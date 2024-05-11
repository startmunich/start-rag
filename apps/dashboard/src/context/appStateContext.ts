import { CrawlerState } from "@/hooks/useCrawlerState";
import { PageCountState } from "@/hooks/usePageCountState";
import { createContext } from "react";

export interface AppState {
    crawler: CrawlerState;
    pageCount: PageCountState;
}

const AppStateContext = createContext<AppState>({
    crawler: {
        isRunning: false,
        inQueue: 0,
        processed: 0,
        cacheMisses: 0,
        lastRunDuration: 0,
        lastRunStartedAt: 0,
        lastRunEndedAt: 0,
        nextRunAt: 0,
    },
    pageCount: {
        count: 0,
    },
});

export default AppStateContext;