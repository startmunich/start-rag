"use client";
import { useCallback, useEffect, useState } from "react";
import debounce from "lodash.debounce";
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

export interface SearchResult {
  hits: SearchHit[];
  limit: number;
  processingTimeMs: number;
  query: string;
}

export interface SearchHit {
  _formatted: FormattedData;
  content: string;
  id: string;
  url: string;
}

export interface FormattedData {
  content: string;
  id: string;
  url: string;
}

const columns: ColumnDef<SearchHit>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: ({ row }) => (
      <div className="w-12 overflow-hidden text-ellipsis [&>em]:bg-yellow-300">
        {row.getValue("id")}
      </div>
    ),
  },
  {
    accessorKey: "content",
    header: "Content",
    cell: ({ row }) => (
      <div
        className="[&>em]:bg-yellow-300"
        dangerouslySetInnerHTML={{ __html: row.getValue("content") }}
      ></div>
    ),
  },
  {
    accessorKey: "url",
    header: "Url",
    cell: ({ row }) => (
      <a
        target="_blank"
        href={row.getValue("url")}
        className="[&>em]:bg-yellow-300"
      >
        {row.getValue("url")}
      </a>
    ),
  },
];

const emptyResult: SearchResult = {
  hits: [],
  limit: 0,
  processingTimeMs: 0,
  query: "",
};

function CrawledPagesTable() {
  const [query, setQuery] = useState("");
  const [searchResult, setSearchResult] = useState<SearchResult>(emptyResult);

  const table = useReactTable({
    data: searchResult.hits,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const updateCrawledPages = useCallback(
    debounce(async (q: string) => {
      setSearchResult(emptyResult);
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!}/search?q=${encodeURIComponent(q)}`,
      );
      const data = await res.json();
      data.hits = data.hits.map((hit: any) => ({
        content: hit._formatted.content,
        id: hit._formatted.id,
        url: hit._formatted.url,
      }));
      setSearchResult(data);
    }, 500),
    [],
  );

  useEffect(() => {
    if (query == "") return;
    updateCrawledPages(query);
  }, [query, updateCrawledPages]);

  return (
    <div className="flex flex-col gap-2 items-end">
      <Input
        onChange={(e) => setQuery(e.target.value)}
        value={query}
        type="search"
        placeholder="Search"
        className="w-48"
      />
      <div className="text-sm text-gray-500">
        {searchResult.hits.length} results in {searchResult.processingTimeMs}ms
      </div>
      <div className="rounded-md border w-full max-h-[70vh] overflow-y-auto">
        <Table className="w-full">
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
          <TableBody className="w-full">
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
                <TableCell colSpan={0} className="h-24 text-center"></TableCell>
                <TableCell colSpan={0} className="h-24 text-center">
                  Search something :)
                </TableCell>
                <TableCell colSpan={0} className="h-24 text-center"></TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

export default CrawledPagesTable;
