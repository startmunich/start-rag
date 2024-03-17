import useCrawlerState, { CrawlerState } from "./useCrawlerState";
import usePageCountState, { PageCountState } from "./usePageCountState";

export interface AppState {
    crawler: CrawlerState;
    pageCount: PageCountState;
}

export default function useAppState(basePath: string): AppState {
    const crawler = useCrawlerState(basePath);
    const pageCount = usePageCountState(basePath);

    return {
        crawler,
        pageCount,
    };
}