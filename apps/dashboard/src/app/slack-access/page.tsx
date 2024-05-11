"use client";
import React from "react";
import UsersTable from "@/components/users-table/users-table";

export default function Crawler() {
  return (
    <main className="h-[100vh]">
      <div className="flex items-center justify-center mt-7">
        <h1 className="text-3xl font-bold text-center text-gray-900">
          Slack Access
        </h1>
      </div>
      <div className="flex justify-center gap-4 mx-7 mt-7">
        <UsersTable />
      </div>
    </main>
  );
}
