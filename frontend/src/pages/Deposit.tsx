import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm, type SubmitHandler } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { depositsApi } from '@/api/endpoints';
import { extractApiError } from '@/api/client';
import { CURRENCIES, type Currency } from '@/types';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';
import {
  depositSchema,
  type DepositInput,
  type DepositValues,
} from '@/validation/depositValidation';

export default function Deposit() {
  const qc = useQueryClient();
  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<DepositInput>({
    resolver: zodResolver(depositSchema),
    defaultValues: { currency: 'USD', amount: undefined as unknown as number },
  });

  const currency = watch('currency');
  const amount = watch('amount');
const idempotencyKey = crypto.randomUUID();
console.log(idempotencyKey);

  const mutation = useMutation({
    mutationFn: (values: DepositValues) =>
      depositsApi.createDeposit(values.currency as Currency, values.amount, idempotencyKey),
    onSuccess: (_d, values) => {
      toast.success('Deposit submitted');
      reset({
        currency: values.currency,
        amount: undefined as unknown as number,
      });
      qc.invalidateQueries({ queryKey: ['wallet'] });
      qc.invalidateQueries({ queryKey: ['transactions'] });
    },
    onError: (err) => toast.error(extractApiError(err, 'Deposit failed')),
  });

  const onSubmit: SubmitHandler<DepositInput> = (data) => {
    const parsed = depositSchema.parse(data);
    mutation.mutate(parsed);
  };

  return (
    <div className="max-w-xl">
      <h1 className="text-2xl font-semibold">Deposit</h1>
      <p className="text-sm text-muted-foreground">
        Add funds to your wallet (simulated)
      </p>

      <form
        onSubmit={handleSubmit(onSubmit)}
        noValidate
        className="mt-6 rounded-xl border border-border bg-card p-6 shadow-sm space-y-4"
      >
        <div>
          <label className="text-sm font-medium" htmlFor="currency">
            Currency
          </label>
          <select
            id="currency"
            {...register('currency')}
            className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          >
            {CURRENCIES.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </select>
          {errors.currency && (
            <p className="mt-1 text-xs text-destructive">
              {errors.currency.message}
            </p>
          )}
        </div>
        <div>
          <label className="text-sm font-medium" htmlFor="amount">
            Amount
          </label>
          <input
            id="amount"
            type="number"
            min="0"
            step="0.01"
            {...register('amount')}
            aria-invalid={!!errors.amount}
            className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
            placeholder="0.00"
          />
          {errors.amount && (
            <p className="mt-1 text-xs text-destructive">
              {errors.amount.message}
            </p>
          )}
        </div>
        <button
          type="submit"
          disabled={mutation.isPending}
          className="w-full inline-flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {mutation.isPending && <Loader2 className="size-4 animate-spin" />}
          Deposit {amount && currency ? `${currency} ${amount}` : ''}
        </button>
      </form>
    </div>
  );
}
