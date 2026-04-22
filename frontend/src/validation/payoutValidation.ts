import { z } from "zod";

export const payoutSchema = z.object({
  sourceCurrency: z.enum(['USD', 'GBP', 'EUR', 'NGN', 'KES'] as const),
  amount: z.coerce
   .number({
    error: (issue) => (issue.input === undefined ? "Amount is required" : "Not a number"),
  })
    .positive("Amount must be greater than 0")
    .max(1_000_000, "Amount is too large"),
  accountNumber: z
    .string()
    .trim()
    .min(6, "Account number must be at least 6 digits")
    .max(20, "Account number must be less than 20 digits")
    .regex(/^[0-9]+$/, "Account number must contain only digits"),
 
});

export type PayoutValues = z.infer<typeof payoutSchema>;
export type PayoutInput = z.input<typeof payoutSchema>;
