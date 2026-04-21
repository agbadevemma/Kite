
import { api } from "./client";
import type {
  AuthResponse,
  Paginated,
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
  getWallet: () => api.get<Wallet>("/wallet").then((r) => r.data),
};


export const transactionsApi = {
  getTransactions: (page = 1, pageSize = 10) =>
    api
      .get<Paginated<Transaction>>("/transactions", { params: { page, pageSize } })
      .then((r) => r.data),
};