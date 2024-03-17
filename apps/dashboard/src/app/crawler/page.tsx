"use client";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import CrawlerContext from "@/context/appStateContext";
import { cn } from "@/lib/utils";
import { PlayCircle, StopCircle } from "lucide-react";
import { useContext } from "react";
import React  from 'react';
import Moment from 'react-moment';

export default function Crawler() {
    const crawlerState = useContext(CrawlerContext);
    const lastRunEndedAt = new Date(crawlerState.lastRunEndedAt);
    const nextRunAt = new Date(crawlerState.nextRunAt);

    return (
      <main className="h-[100vh]">
        <div className="flex items-center justify-center mt-7">
          <h1 className="text-3xl font-bold text-center text-gray-900">
            Crawler
          </h1>
        </div>
        <div className="flex justify-center gap-4 mx-7 mt-7">
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>Current Status</CardTitle>
              <CardDescription className={cn(
                "flex items-center gap-1",
                crawlerState.isRunning ? "text-green-800" : "text-red-800",
              )}>
                {crawlerState.isRunning ? <>
                  <PlayCircle className="w-5 h-5 inline-block" />
                  Running
                </> : <>
                  <StopCircle className="w-5 h-5 inline-block" />
                  Stopped
                </>}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2">
                <div className="flex-grow">
                  <h2 className="mb-1">Queue</h2>
                  <div className="flex justify-center gap-2 text-2xl px-3 py-2 bg-black/10 rounded-xl">
                    {crawlerState.inQueue}
                  </div>
                </div>
                <div className="flex-grow">
                  <h2 className="mb-1">Done</h2>
                  <div className="flex justify-center gap-2 text-2xl px-3 py-2 bg-black/10 rounded-xl">
                    {crawlerState.processed}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>Schedule</CardTitle>
              <CardDescription>
                {
                  crawlerState.isRunning
                  ? "Crawler is currently running"
                  : <>
                    Next run in&nbsp;
                    <Moment fromNow date={nextRunAt} />
                  </>
                }
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex gap-7">
                <div className="flex flex-col items-start">
                  <div className="px-3 py-2">
                    Last run ended&nbsp;
                    <b><Moment fromNow date={lastRunEndedAt} /></b>
                    &nbsp;and took <b>{crawlerState.lastRunDuration / 1000 / 60} min</b>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
        
      </main>
    );
  }