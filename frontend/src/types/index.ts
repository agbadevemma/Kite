export type Currency = "USD" | "GBP" | "EUR" | "NGN" | "KES";
export const CURRENCIES: Currency[] = ["USD", "GBP", "EUR", "NGN", "KES"];

export interface User {
  id: string;
  email: string;
  name?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Balance {
  currency: Currency;
  amount: number;
}

export interface Wallet {
  balances: Balance[];
}

export interface Quote {
  id: string;
  from_currency: Currency;
  to_currency: Currency;
  amount_in: number;
  amount_out: number;
  rate: number;
  fee: number;
  expires_at: string; // ISO
}

export type TxType = "deposit" | "conversion" | "payout";
export type TxStatus = "pending" | "completed" | "failed";

export interface Transaction {
  id: string;
  type: TxType;
  status: TxStatus;
  amount: number;
  currency: Currency;
  created_at: string;
  meta?: Record<string, unknown>;
}

export interface Paginated<T> {
  data: T[];
  page: number;
  pageSize: number;
  total: number;
}

export interface ApiError {
  code?: string;
  message: string;
}