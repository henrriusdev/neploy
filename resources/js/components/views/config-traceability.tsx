import React from "react";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from "@/components/ui/table";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button";
import {Trace, TracesSettingsProps} from "@/types";
import {
  ColumnDef, ColumnFiltersState,
  flexRender,
  getCoreRowModel, getFilteredRowModel,
  getPaginationRowModel, getSortedRowModel,
  SortingState,
  useReactTable
} from "@tanstack/react-table";
import {useTranslation} from "react-i18next";
import {ArrowUpDown, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight} from "lucide-react";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";

const TraceabilityTab = ({traces}: TracesSettingsProps) => {
  const {t} = useTranslation()
  const columns: ColumnDef<Trace>[] = [
    {
      accessorKey: "actionTimestamp",
      header: ({column}) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            {t("dashboard.settings.trace.date")}
            <ArrowUpDown className="ml-2 h-4 w-4"/>
          </Button>
        )
      },
    },
    {
      accessorKey: "email",
      header: ({column}) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            {t("dashboard.settings.trace.user")}
            <ArrowUpDown className="ml-2 h-4 w-4"/>
          </Button>
        )
      },
    },
    {
      accessorKey: "action",
      header: ({column}) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            {t("dashboard.settings.trace.action")}
            <ArrowUpDown className="ml-2 h-4 w-4"/>
          </Button>
        )
      },
    },
    {
      accessorKey: "sqlStatement",
      header: t("dashboard.settings.trace.query"),
    }
  ]

  interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[]
    data: TData[]
  }

  function DataTable<TData, TValue>({
                                      columns,
                                      data
                                    }: DataTableProps<TData, TValue>) {
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
      []
    )
    const table = useReactTable({
      data,
      columns,
      getCoreRowModel: getCoreRowModel(),
      getPaginationRowModel: getPaginationRowModel(),
      onSortingChange: setSorting,
      getSortedRowModel: getSortedRowModel(),
      onColumnFiltersChange: setColumnFilters,
      getFilteredRowModel: getFilteredRowModel(),
      state: {
        sorting,
        columnFilters,
      },
    })

    return (
      <div>
        <div className="flex items-center py-4">
          <Input
            placeholder="Filter queries..."
            value={(table.getColumn("sqlStatement")?.getFilterValue() as string) ?? ""}
            onChange={(event) =>
              table.getColumn("sqlStatement")?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
        </div>
        <div className="rounded-md border">
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
                            header.getContext()
                          )}
                      </TableHead>
                    )
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
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={columns.length} className="h-24 text-center">
                    No results.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
        <div className="flex items-center justify-end p-4">
          <div className="flex items-center space-x-6 lg:space-x-8">
            <div className="flex items-center space-x-2">
              <p className="text-sm font-medium">Rows per page</p>
              <Select
                value={`${table.getState().pagination.pageSize}`}
                onValueChange={(value) => {
                  table.setPageSize(Number(value))
                }}
              >
                <SelectTrigger className="h-8 w-[70px]">
                  <SelectValue placeholder={table.getState().pagination.pageSize}/>
                </SelectTrigger>
                <SelectContent side="top">
                  {[10, 20, 30, 40, 50].map((pageSize) => (
                    <SelectItem key={pageSize} value={`${pageSize}`}>
                      {pageSize}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="flex w-[100px] items-center justify-center text-sm font-medium">
              Page {table.getState().pagination.pageIndex + 1} of{" "}
              {table.getPageCount()}
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                className="hidden h-8 w-8 p-0 lg:flex"
                onClick={() => table.setPageIndex(0)}
                disabled={!table.getCanPreviousPage()}
              >
                <span className="sr-only">Go to first page</span>
                <ChevronsLeft/>
              </Button>
              <Button
                variant="outline"
                className="h-8 w-8 p-0"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
              >
                <span className="sr-only">Go to previous page</span>
                <ChevronLeft/>
              </Button>
              <Button
                variant="outline"
                className="h-8 w-8 p-0"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
              >
                <span className="sr-only">Go to next page</span>
                <ChevronRight/>
              </Button>
              <Button
                variant="outline"
                className="hidden h-8 w-8 p-0 lg:flex"
                onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                disabled={!table.getCanNextPage()}
              >
                <span className="sr-only">Go to last page</span>
                <ChevronsRight/>
              </Button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between items-center">
            {t("dashboard.settings.trace.title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable columns={columns} data={traces}/>
        </CardContent>
      </Card>
    </div>
  );
};

export default TraceabilityTab;
