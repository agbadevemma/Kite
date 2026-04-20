import { TOKEN_KEY } from "@/api/client";
import { useCallback, useEffect, useState } from "react";


export function useAuth() {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem(TOKEN_KEY)
  );

  useEffect(() => {
    const onStorage = () => setToken(localStorage.getItem(TOKEN_KEY));
    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

  const setAuthToken = useCallback((t: string | null) => {
    if (t) localStorage.setItem(TOKEN_KEY, t);
    else localStorage.removeItem(TOKEN_KEY);
    setToken(t);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);
    window.location.href = "/login";
  }, []);

  return { token, isAuthenticated: !!token, setAuthToken, logout };
}