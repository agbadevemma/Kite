import { useState } from "react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { transactionsApi } from "@/api/endpoints";
import { extractApiError } from "@/api/client";
import { formatDate, formatMoney } from "@/utils/format";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { ErrorState } from "@/components/Errorstate";

const PAGE_SIZE = 10;

export default function Transactions() {
  const [page, setPage] = useState(1);
  const query = useQuery({
    queryKey: ["transactions", { page, pageSize: PAGE_SIZE }],
    queryFn: () => transactionsApi.getTransactions(page, PAGE_SIZE),
    placeholderData: keepPreviousData,
  });

  const totalPages = query.data ? Math.max(1, Math.ceil(query.data.total / PAGE_SIZE)) : 1;

  return (
    <div>
      <h1 className="text-2xl font-semibold">Transactions</h1>
      <p className="text-sm text-muted-foreground">All your wallet activity</p>

      <div className="mt-6 rounded-xl border border-border bg-card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-muted/50 text-muted-foreground text-xs uppercase">
              <tr>
                <th className="text-left font-medium px-4 py-3">Type</th>
                <th className="text-left font-medium px-4 py-3">Status</th>
                <th className="text-right font-medium px-4 py-3">Amount</th>
                <th className="text-left font-medium px-4 py-3">Currency</th>
                <th className="text-left font-medium px-4 py-3">Date</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {query.isLoading &&
                [...Array(PAGE_SIZE)].map((_, i) => (
                  <tr key={i}>
                    <td colSpan={5} className="px-4 py-3">
                      <Skeleton className="h-5" />
                    </td>
                  </tr>
                ))}
              {query.data &&
                query.data.data.map((t) => (
                  <tr key={t.id} className="hover:bg-muted/30">
                    <td className="px-4 py-3 capitalize font-medium">{t.type}</td>
                    <td className="px-4 py-3">
                      <StatusBadge status={t.status} />
                    </td>
                    <td className="px-4 py-3 text-right font-semibold">
                      {formatMoney(t.amount, t.currency)}
                    </td>
                    <td className="px-4 py-3 text-muted-foreground">{t.currency}</td>
                    <td className="px-4 py-3 text-muted-foreground">{formatDate(t.created_at)}</td>
                  </tr>
                ))}
              {query.data && query.data.data.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-4 py-12 text-center text-muted-foreground">
                    No transactions yet
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {query.isError && (
          <div className="p-4">
            <ErrorState message={extractApiError(query.error, "Failed to load transactions")} />
          </div>
        )}

        <div className="flex items-center justify-between gap-3 border-t border-border px-4 py-3">
          <div className="text-xs text-muted-foreground">
            Page {page} of {totalPages}
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1 || query.isFetching}
              className="inline-flex items-center gap-1 rounded-md border border-input bg-card px-3 py-1.5 text-xs font-medium hover:bg-muted disabled:opacity-50"
            >
              <ChevronLeft className="size-3.5" /> Prev
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages || query.isFetching}
              className="inline-flex items-center gap-1 rounded-md border border-input bg-card px-3 py-1.5 text-xs font-medium hover:bg-muted disabled:opacity-50"
            >
              Next <ChevronRight className="size-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const map: Record<string, string> = {
    completed: "bg-primary/10 text-primary",
    pending: "bg-warning/10 text-[hsl(var(--warning))]",
    failed: "bg-destructive/10 text-destructive",
  };
  return (
    <span className={`inline-block text-[10px] px-2 py-0.5 rounded-full font-medium capitalize ${map[status] ?? "bg-muted text-muted-foreground"}`}>
      {status}
    </span>
  );
}