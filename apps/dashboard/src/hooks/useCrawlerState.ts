import { useEffect, useState } from "react";

export interface CrawlerState {
    isRunning: boolean;
    inQueue: number;
    processed: number;
    lastRunDuration: number;
    lastRunStartedAt: number;
    lastRunEndedAt: number;
    nextRunAt: number;
};  

export default function useCrawlerState(basePath: string): CrawlerState {
    const [state, setState] = useState<CrawlerState>({
        isRunning: false,
        inQueue: 0,
        processed: 0,
        lastRunDuration: 0,
        lastRunStartedAt: 0,
        lastRunEndedAt: 0,
        nextRunAt: 0,
    });

    useEffect(() => {
        const timeout = setTimeout(async () => {
            const result = await fetch(`${basePath}/state`);
            const newState = await (result.json() as Promise<CrawlerState>);
            setState(newState);
        }, 3000);

        return () => clearTimeout(timeout);
    }, [basePath]);

    return state;
}