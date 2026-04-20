import type { Currency } from "../types";

const localeMap: Record<Currency, string> = {
  USD: "en-US",
  GBP: "en-GB",
  EUR: "de-DE",
  NGN: "en-NG",
  KES: "en-KE",
};

export function formatMoney(amount: number, currency: Currency): string {
  try {
    return new Intl.NumberFormat(localeMap[currency] ?? "en-US", {
      style: "currency",
      currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    return `${currency} ${amount.toFixed(2)}`;
  }
}

export function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return iso;
  }
}

export const currencyFlag: Record<Currency, string> = {
  USD: "🇺🇸",
  GBP: "🇬🇧",
  EUR: "🇪🇺",
  NGN: "🇳🇬",
  KES: "🇰🇪",
};