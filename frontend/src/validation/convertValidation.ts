
import { CURRENCIES } from '@/types';
import { z } from 'zod';

    

export const convertSchema = z
  .object({
    source: z.enum(CURRENCIES),
    target: z.enum(CURRENCIES),
    amount: z.coerce
      .number({
        error: (issue) =>
          issue.input === undefined ? 'Amount is required' : 'Not a number',
      })
      .positive('Amount must be greater than 0')
      .max(1_000_000, 'Amount is too large'),
  })
  .refine((d) => d.source !== d.target, {
    message: 'Source and target must differ',
    path: ['target'],
  });

export type ConvertValues = z.infer<typeof convertSchema>;
export type ConvertInput = z.input<typeof convertSchema>;
