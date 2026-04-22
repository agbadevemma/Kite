import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { transactionsApi, walletApi } from "@/api/endpoints";
import { CURRENCIES } from "@/types";
import { currencyFlag, formatDate, formatMoney } from "@/utils/format";
import { extractApiError } from "@/api/client";
import { ArrowLeftRight, PlusCircle, Send } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { ErrorState } from "@/components/Errorstate";

export default function Dashboard() {
  const wallet = useQuery({
    queryKey: ["wallet"],
    queryFn: walletApi.getWallet,
  });

  const txs = useQuery({
    queryKey: ["transactions", { page: 1, pageSize: 5 }],
    queryFn: () => transactionsApi.getTransactions(1, 5),
  });

  const balanceMap = new Map(wallet.data?.balances.map((b) => [b.currency, b.amount]));
console.log(txs?.data?.data);

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">Dashboard</h1>
          <p className="text-sm text-muted-foreground">Your balances and recent activity</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Link
            to="/deposit"
            className="inline-flex items-center gap-2 rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            <PlusCircle className="size-4" /> Deposit
          </Link>
          <Link
            to="/convert"
            className="inline-flex items-center gap-2 rounded-md border border-input bg-card px-3 py-2 text-sm font-medium hover:bg-muted"
          >
            <ArrowLeftRight className="size-4" /> Convert
          </Link>
          <Link
            to="/payout"
            className="inline-flex items-center gap-2 rounded-md border border-input bg-card px-3 py-2 text-sm font-medium hover:bg-muted"
          >
            <Send className="size-4" /> Payout
          </Link>
        </div>
      </div>

      {wallet.isError && <ErrorState message={extractApiError(wallet.error, "Failed to load wallet")} />}

      <section>
        <h2 className="text-sm font-medium text-muted-foreground mb-3">Balances</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
          {wallet.isLoading
            ? CURRENCIES.map((c) => <Skeleton key={c} className="h-28" />)
            : CURRENCIES.map((c) => {
                const amt = balanceMap.get(c) ?? 0;
                return (
                  <div
                    key={c}
                    className="rounded-xl border border-border bg-card p-5 shadow-sm hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-2xl">{currencyFlag[c]}</span>
                      <span className="text-xs font-medium text-muted-foreground">{c}</span>
                    </div>
                    <div className="mt-3 text-2xl font-semibold tracking-tight">
                      {formatMoney(amt, c)}
                    </div>
                    <div className="text-xs text-muted-foreground mt-1">Available balance</div>
                  </div>
                );
              })}
        </div>
      </section>

      <section>
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-medium text-muted-foreground">Recent transactions</h2>
          <Link to="/transactions" className="text-sm text-primary hover:underline">
            View all
          </Link>
        </div>
        <div className="rounded-xl border border-border bg-card overflow-hidden">
          {txs.isLoading && (
            <div className="p-4 space-y-3">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-10" />
              ))}
            </div>
          )}
          {txs.isError && (
            <div className="p-4">
              <ErrorState message={extractApiError(txs.error, "Failed to load transactions")} />
            </div>
          )}
          {txs.data && txs.data.data?.length === 0 && (
            <div className="p-8 text-center text-sm text-muted-foreground">
              No transactions yet
            </div>
          )}
          {txs?.data && txs.data?.data?.length > 0 && (
            <ul className="divide-y divide-border">
              {txs.data.data.map((t) => (
                <li key={t.id} className="flex items-center justify-between gap-4 px-4 py-3">
                  <div className="flex items-center gap-3 min-w-0">
                    <span className="grid size-9 place-items-center rounded-full bg-accent text-accent-foreground capitalize text-xs font-medium">
                      {t.type[0]}
                    </span>
                    <div className="min-w-0">
                      <div className="text-sm font-medium capitalize truncate">{t.type}</div>
                      <div className="text-xs text-muted-foreground">{formatDate(t.created_at)}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-semibold">{formatMoney(t.amount, t.currency)}</div>
                    <StatusBadge status={t.status} />
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </section>
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
