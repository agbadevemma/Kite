import * as z from 'zod';
export const depositSchema = z.object({
  currency: z.enum(['USD', 'GBP', 'EUR', 'NGN', 'KES'] as const),
  amount: z.coerce
    .number()
    .positive('Amount must be greater than 0')
    .max(1_000_000, 'Amount is too large'),
});

export type DepositInput = z.input<typeof depositSchema>;
export type DepositValues = z.output<typeof depositSchema>;
