import { useEffect, useState } from "react";

export function useCountdown(targetIso?: string) {
  const target = targetIso ? new Date(targetIso).getTime() : 0;
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!targetIso) return;
    const id = setInterval(() => setNow(Date.now()), 250);
    return () => clearInterval(id);
  }, [targetIso]);

  const remainingMs = Math.max(0, target - now);
  const seconds = Math.floor(remainingMs / 1000);
  return {
    remainingMs,
    seconds,
    expired: !!targetIso && remainingMs <= 0,
    label: `${String(Math.floor(seconds / 60)).padStart(2, "0")}:${String(seconds % 60).padStart(2, "0")}`,
  };
}