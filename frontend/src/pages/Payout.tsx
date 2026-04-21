import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { payoutsApi } from "@/api/endpoints";
import { extractApiError } from "@/api/client";
import { CURRENCIES, type Currency } from "@/types";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { payoutSchema, type PayoutInput, type PayoutValues } from "@/validation/payoutValidation";



export default function Payout() {
  const qc = useQueryClient();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<PayoutInput>({
    resolver: zodResolver(payoutSchema),
    defaultValues: {
      sourceCurrency: "USD",
      amount: undefined as unknown as number,
      accountName: "",
      accountNumber: "",
      bankCode: "",
    },
  });

  const mutation = useMutation({
    mutationFn: (values: PayoutValues) =>
      payoutsApi.createPayout({
        sourceCurrency: values.sourceCurrency as Currency,
        amount: values.amount,
        accountNumber: values.accountNumber,
        bankCode: values.bankCode,
        accountName: values.accountName,
      }),
    onSuccess: (_d, values) => {
      toast.success("Payout initiated");
      reset({
        sourceCurrency: values.sourceCurrency,
        amount: undefined as unknown as number,
        accountName: "",
        accountNumber: "",
        bankCode: "",
      });
      qc.invalidateQueries({ queryKey: ["wallet"] });
      qc.invalidateQueries({ queryKey: ["transactions"] });
    },
    onError: (err) => {
      const msg = extractApiError(err, "Payout failed");
      if (msg.toLowerCase().includes("insufficient")) {
        toast.error("Insufficient balance for this payout.");
      } else {
        toast.error(msg);
      }
    },
  });
    const onSubmit: SubmitHandler<PayoutInput> = (data) => {
      const parsed = payoutSchema.parse(data);
      mutation.mutate(parsed);
    };

  return (
    <div className="max-w-xl">
      <h1 className="text-2xl font-semibold">Send payout</h1>
      <p className="text-sm text-muted-foreground">Transfer to a bank account</p>

      <form
        onSubmit={handleSubmit(onSubmit)}
        noValidate
        className="mt-6 rounded-xl border border-border bg-card p-6 shadow-sm space-y-4"
      >
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="text-sm font-medium" htmlFor="sourceCurrency">Source currency</label>
            <select
              id="sourceCurrency"
              {...register("sourceCurrency")}
              className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              {CURRENCIES.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
            {errors.sourceCurrency && (
              <p className="mt-1 text-xs text-destructive">{errors.sourceCurrency.message}</p>
            )}
          </div>
          <div>
            <label className="text-sm font-medium" htmlFor="amount">Amount</label>
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
        </div>
        <div>
          <label className="text-sm font-medium" htmlFor="accountName">Account name</label>
          <input
            id="accountName"
            {...register("accountName")}
            aria-invalid={!!errors.accountName}
            className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            placeholder="Jane Doe"
          />
          {errors.accountName && (
            <p className="mt-1 text-xs text-destructive">{errors.accountName.message}</p>
          )}
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="text-sm font-medium" htmlFor="accountNumber">Account number</label>
            <input
              id="accountNumber"
              inputMode="numeric"
              {...register("accountNumber")}
              aria-invalid={!!errors.accountNumber}
              className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              placeholder="0123456789"
            />
            {errors.accountNumber && (
              <p className="mt-1 text-xs text-destructive">{errors.accountNumber.message}</p>
            )}
          </div>
          <div>
            <label className="text-sm font-medium" htmlFor="bankCode">Bank code</label>
            <input
              id="bankCode"
              {...register("bankCode")}
              aria-invalid={!!errors.bankCode}
              className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              placeholder="044"
            />
            {errors.bankCode && (
              <p className="mt-1 text-xs text-destructive">{errors.bankCode.message}</p>
            )}
          </div>
        </div>
        <button
          type="submit"
          disabled={mutation.isPending}
          className="w-full inline-flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {mutation.isPending && <Loader2 className="size-4 animate-spin" />}
          Send payout
        </button>
      </form>
    </div>
  );
}