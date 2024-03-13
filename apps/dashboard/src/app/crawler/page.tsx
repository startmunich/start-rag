"use client";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import CrawlerContext from "@/context/crawlerContext";
import { useContext } from "react";

export default function Crawler() {
    const crawlerState = useContext(CrawlerContext);

    return (
      <main className="h-[100vh]">
        <div className="flex items-center justify-center mt-7">
          <h1 className="text-3xl font-bold text-center text-gray-900">
            Crawler Status
          </h1>
        </div>
        <div className="flex justify-start gap-4 mx-7 mt-7">
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>Current Progress</CardTitle>
              <CardDescription>{crawlerState.isRunning ? "Running" : "Stopped"}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex justify-center gap-2 text-2xl px-3 py-2 bg-black/10 rounded-xl">
                <div className="">{crawlerState.processed}</div>
                <div className="">/</div>
                <div className="">{crawlerState.inQueue}</div>
              </div>
            </CardContent>
          </Card>
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>Last Run</CardTitle>
              <CardDescription></CardDescription>
            </CardHeader>
            <CardContent>
              
            </CardContent>
          </Card>
        </div>
      </main>
    );
  }
  