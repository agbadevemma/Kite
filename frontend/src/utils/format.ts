import type { Currency } from "../types";

const localeMap: Record<Currency, string> = {
  USD: "en-US",
  GBP: "en-GB",
  EUR: "en-US",
  NGN: "en-NG",
  KES: "en-KE",
};

const currencyDecimals: Record<Currency, number> = {
  USD: 2,
  GBP: 2,
  EUR: 2,
  NGN: 2,
  KES: 2,
};

export function formatMoney(
  amount: number,
  currency: Currency
): string {
  try {
    const decimals = currencyDecimals[currency] ?? 2;

    const formattedAmount =
      amount / Math.pow(10, decimals);

    return new Intl.NumberFormat(
      localeMap[currency] ?? "en-US",
      {
        style: "currency",
        currency,
        minimumFractionDigits: decimals,
        maximumFractionDigits: decimals,
      }
    ).format(formattedAmount);
  } catch {
    return `${currency} ${amount}`;
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