import { useEffect, useState } from "react";

export interface PageCountState {
    count: number;
};  

export default function usePageCountState(basePath: string): PageCountState {
    const [state, setState] = useState<PageCountState>({
        count: 0,
    });

    useEffect(() => {
        const updateState = async () => {
            const result = await fetch(`${basePath}/pages/count`);
            const newState = await (result.json() as Promise<PageCountState>);
            setState(newState);
        };
        const timeout = setInterval(updateState, 3000);
        updateState();

        return () => clearInterval(timeout);
    }, [basePath]);

    return state;
}