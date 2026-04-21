
import { api } from "./client";
import type {
  AuthResponse,
  Currency,
  Paginated,
  Quote,
  Transaction,
  Wallet,
} from "@/types";

export const authApi = {
  login: (email: string, password: string) =>
    api.post<AuthResponse>("/auth/login", { email, password }).then((r) => r.data),
  signup: (email: string, password: string, name?: string) =>
    api.post<AuthResponse>("/auth/signup", { email, password, name }).then((r) => r.data),
};


export const walletApi = {
  getWallet: () => api.get<Wallet>("/wallet/balances").then((r) => r.data),
};


export const depositsApi = {
  createDeposit: (currency: Currency, amount: number) =>
    api.post<Transaction>("/deposits", { currency, amount }).then((r) => r.data),
};

export const conversionsApi = {
  getQuote: (sourceCurrency: Currency, targetCurrency: Currency, amount: number) =>
    api
      .post<Quote>("/conversions/quote", { sourceCurrency, targetCurrency, amount })
      .then((r) => r.data),
  executeConversion: (quoteId: string) =>
    api.post<Transaction>("/conversions/execute", { quoteId }).then((r) => r.data),
};

export const payoutsApi = {
  createPayout: (params: {
    sourceCurrency: Currency;
    amount: number;
    accountNumber: string;
    bankCode: string;
    accountName: string;
  }) => api.post<Transaction>("/payouts", params).then((r) => r.data),
};

  
export const transactionsApi = {
  getTransactions: (page = 1, pageSize = 10) =>
    api
      .get<Paginated<Transaction>>("/transactions", { params: { page, pageSize } })
      .then((r) => r.data),
};