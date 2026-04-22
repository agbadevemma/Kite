import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { conversionsApi } from "@/api/endpoints";
import { extractApiError } from "@/api/client";
import { CURRENCIES, type Currency, type Quote } from "@/types";
import { formatMoney } from "@/utils/format";
import { toast } from "sonner";
import { Loader2, ArrowRight, Timer } from "lucide-react";
import { convertSchema,type ConvertInput, type ConvertValues } from "@/validation/convertValidation";
import { useCountdown } from "@/hooks/useCountDown";

type Step = "input" | "quote";

export default function Convert() {
  const [step, setStep] = useState<Step>("input");
  const [quote, setQuote] = useState<Quote | null>(null);

  const qc = useQueryClient();
  const countdown = useCountdown(quote?.expires_at);

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<ConvertInput>({
    resolver: zodResolver(convertSchema),
    defaultValues: {
      source: "USD",
      target: "EUR",
      amount: undefined as unknown as number,
    },
  });
  
  const source = watch("source");

  const quoteMutation = useMutation({
    mutationFn: (values: ConvertValues) =>
      conversionsApi.getQuote(values.source as Currency, values.target as Currency, values.amount),
    onSuccess: (q) => {
      setQuote(q);
      setStep("quote");
    },
    onError: (err) => toast.error(extractApiError(err, "Could not fetch quote")),
  });


  
  const executeMutation = useMutation({
    mutationFn: () => conversionsApi.executeConversion(quote!.id),
    onSuccess: () => {
      toast.success("Conversion completed");
      qc.invalidateQueries({ queryKey: ["wallet"] });
      qc.invalidateQueries({ queryKey: ["transactions"] });
      setStep("input");
      setQuote(null);
      reset({ source: "USD", target: "EUR", amount: undefined as unknown as number });
    },
    onError: (err) => {
      const msg = extractApiError(err, "Conversion failed");
      const lower = msg.toLowerCase();
      if (lower.includes("expired")) {
        toast.error("Quote expired. Please request a new one.");
        setStep("input");
        setQuote(null);
      } else if (lower.includes("insufficient")) {
        toast.error("Insufficient balance for this conversion.");
      } else {
        toast.error(msg);
      }
    },
  });

  const onSubmit = (data: ConvertInput) => {
    const parsed = convertSchema.parse(data);
    quoteMutation.mutate(parsed);
  }
  return (
    <div className="max-w-xl">
      <h1 className="text-2xl font-semibold">Convert currency</h1>
      <p className="text-sm text-muted-foreground">Get a live FX quote and confirm</p>

      {step === "input" && (
        <form
          onSubmit={handleSubmit(onSubmit)}
          noValidate
          className="mt-6 rounded-xl border border-border bg-card p-6 shadow-sm space-y-4"
        >
          <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto_1fr] gap-3 items-end">
            <div>
              <label className="text-sm font-medium" htmlFor="source">From</label>
              <select
                id="source"
                {...register("source")}
                className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                {CURRENCIES.map((c) => (
                  <option key={c} value={c}>
                    {c}
                  </option>
                ))}
              </select>
            </div>
            <div className="hidden sm:flex items-center justify-center pb-2">
              <ArrowRight className="size-4 text-muted-foreground" />
            </div>
            <div>
              <label className="text-sm font-medium" htmlFor="target">To</label>
              <select
                id="target"
                {...register("target")}
                className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                {CURRENCIES.map((c) => (
                  <option key={c} value={c}>
                    {c}
                  </option>
                ))}
              </select>
            </div>
          </div>
          {errors.target && (
            <p className="text-xs text-destructive">{errors.target.message}</p>
          )}
          <div>
            <label className="text-sm font-medium" htmlFor="amount">Amount ({source})</label>
            <input
              id="amount"
              type="number"
              min="0"
              step="0.01"
              {...register("amount")}
              aria-invalid={!!errors.amount}
              className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              placeholder="0.00"
            />
            {errors.amount && (
              <p className="mt-1 text-xs text-destructive">{errors.amount.message}</p>
            )}
          </div>
          <button
            type="submit"
            disabled={quoteMutation.isPending}
            className="w-full inline-flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
          >
            {quoteMutation.isPending && <Loader2 className="size-4 animate-spin" />}
            Get quote
          </button>
        </form>
      )}

      {step === "quote" && quote && (
        <div className="mt-6 rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Review conversion</h2>
            <div
              className={`inline-flex items-center gap-1.5 text-xs font-medium px-2 py-1 rounded-full ${
                countdown.expired
                  ? "bg-destructive/10 text-destructive"
                  : "bg-accent text-accent-foreground"
              }`}
            >
              <Timer className="size-3.5" />
              {countdown.expired ? "Expired" : `Expires in ${countdown.label}`}
            </div>
          </div>

          <div className="rounded-lg bg-muted/50 p-4 space-y-2">
            <Row label="You send" value={formatMoney(quote.amount_in, quote.from_currency)} />
            <Row label="You receive" value={formatMoney(quote.amount_out, quote.to_currency)} highlight />
            <Row
              label="Rate"
              value={`1 ${quote.from_currency} = ${quote.rate} ${quote.to_currency}`}
            />
            <Row label="Fee" value={formatMoney(quote.fee, quote.from_currency)} />
          </div>

          <div className="flex flex-col sm:flex-row gap-2">
            <button
              type="button"
              onClick={() => {
                setStep("input");
                setQuote(null);
              }}
              className="flex-1 rounded-md border border-input bg-card px-4 py-2.5 text-sm font-medium hover:bg-muted"
            >
              Cancel
            </button>
            <button
              type="button"
              disabled={countdown.expired || executeMutation.isPending}
              onClick={() => executeMutation.mutate()}
              className="flex-1 inline-flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed"
            >
              {executeMutation.isPending && <Loader2 className="size-4 animate-spin" />}
              Confirm conversion
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

function Row({ label, value, highlight }: { label: string; value: string; highlight?: boolean }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className={highlight ? "text-base font-semibold" : "text-sm font-medium"}>{value}</span>
    </div>
  );
}