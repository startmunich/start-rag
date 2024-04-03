"use client";
import {
  AlertDialogFooter,
  AlertDialogHeader,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import CrawlerContext from "@/context/appStateContext";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useContext, useState } from "react";
import React from "react";

export default function Crawler() {
  const [loading, setLoading] = useState<boolean>(false);

  const purgeDb = async () => {
    setLoading(true);
    try {
      await fetch(
        `${process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!}/db/purge`,
        {
          method: "POST",
        },
      );
    } catch (error) {
      console.error(error);
    }
    setLoading(false);
  };

  return (
    <main className="h-[100vh]">
      <div className="flex items-center justify-center mt-7">
        <h1 className="text-3xl font-bold text-center text-gray-900">Manage</h1>
      </div>
      <div className="flex justify-center gap-4 mx-7 mt-7">
        <Card className="w-[350px] relative">
          <CardHeader>
            <CardTitle>Purge Database</CardTitle>
            <CardDescription>
              Purge the database of all data. This action is irreversible.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form>
              <div className="grid w-full items-center gap-4"></div>
            </form>
          </CardContent>
          <CardFooter className="flex justify-between">
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button disabled={loading} variant="destructive">
                  {loading ? "Purging..." : "Purge Database"}
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This action is irreversible. Once you purge the database,
                    you will have to crawl the workspace again.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction onClick={purgeDb}>
                    Continue
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
