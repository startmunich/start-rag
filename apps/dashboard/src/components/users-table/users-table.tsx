"use client";
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
    header: "SlackID",
    cell: ({ row }) => (
      <div className="w-12 overflow-hidden text-ellipsis [&>em]:bg-yellow-300">
        {row.getValue("id")}
      </div>
    ),
  },
  {
    accessorKey: "username",
    header: "Username",
    cell: ({ row }) => (
      <div
        className="[&>em]:bg-yellow-300"
        dangerouslySetInnerHTML={{ __html: row.getValue("content") }}
      ></div>
    ),
  },
];

function UsersTable() {
  const table = useReactTable({
    data: [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
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
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={0} className="h-24 text-center"></TableCell>
              <TableCell colSpan={0} className="h-24 text-center">
                Awesome People :)
              </TableCell>
              <TableCell colSpan={0} className="h-24 text-center"></TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export default UsersTable;
