
import { api } from "./client";
import type {
  AuthResponse,
} from "@/types";

export const authApi = {
  login: (email: string, password: string) =>
    api.post<AuthResponse>("/auth/login", { email, password }).then((r) => r.data),
  signup: (email: string, password: string, name?: string) =>
    api.post<AuthResponse>("/auth/signup", { email, password, name }).then((r) => r.data),
};