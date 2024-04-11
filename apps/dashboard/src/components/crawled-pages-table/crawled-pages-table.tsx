"use client";
import { useEffect, useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import "reactflow/dist/style.css";
import { Input } from "../ui/input";

export interface CrawledPage {
  child_pages: string;
  page_id: string;
  url: string;
}

const columns: ColumnDef<CrawledPage>[] = [
  {
    accessorKey: "page_id",
    header: "Page ID",
    cell: ({ row }) => <div>{row.getValue("page_id")}</div>,
  },
  {
    accessorKey: "child_pages",
    header: "Child Pages",
    cell: ({ row }) => <div>{row.getValue("child_pages")}</div>,
  },
  {
    accessorKey: "url",
    header: "Url",
    cell: ({ row }) => <div>{row.getValue("url")}</div>,
  },
];

function CrawledPagesTable() {
  const [query, setQuery] = useState("");
  const [crawledPages, setCrawledPages] = useState<CrawledPage[]>([]);

  const updateCrawledPages = async () => {
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!}/pages`,
    );
    const data = await res.json();
    setCrawledPages(data.pages);
  };

  useEffect(() => {
    updateCrawledPages();
  }, []);

  const filteredPages = crawledPages.filter((page) =>
    [page.page_id, page.url].some((field) =>
      field.toLowerCase().includes(query.toLowerCase()),
    ),
  );

  const table = useReactTable({
    data: filteredPages,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="flex flex-col gap-2 items-end">
      <Input
        onChange={(e) => setQuery(e.target.value)}
        value={query}
        type="search"
        placeholder="Search"
        className="w-48"
      />
      <div className="rounded-md border w-full">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

export default CrawledPagesTable;
