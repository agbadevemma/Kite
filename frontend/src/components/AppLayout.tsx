import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/utils";
import { Wallet, ArrowLeftRight, Send, ListOrdered, PlusCircle, LogOut } from "lucide-react";

const nav = [
  { to: "/", label: "Dashboard", icon: Wallet },
  { to: "/deposit", label: "Deposit", icon: PlusCircle },
  { to: "/convert", label: "Convert", icon: ArrowLeftRight },
  { to: "/payout", label: "Payout", icon: Send },
  { to: "/transactions", label: "Transactions", icon: ListOrdered },
];

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const { logout } = useAuth();
  const { pathname } = useLocation();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-background text-foreground w-full">
      <header className="border-b border-border bg-card sticky top-0 z-30 ">
        <div className="container mx-auto flex h-16 items-center justify-between gap-4">
          <button
            onClick={() => navigate("/")}
            className="flex items-center gap-2 font-semibold text-lg"
          >
            <span className="grid size-8 place-items-center rounded-md bg-primary text-primary-foreground">
              ₲
            </span>
            <span>Kite</span>
          </button>
          <nav className="hidden md:flex items-center gap-1">
            {nav.map((n) => {
              const active = pathname === n.to;
              const Icon = n.icon;
              return (
                <Link
                  key={n.to}
                  to={n.to}
                  className={cn(
                    "flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                    active
                      ? "bg-accent text-accent-foreground"
                      : "text-muted-foreground hover:text-foreground hover:bg-muted"
                  )}
                >
                  <Icon className="size-4" />
                  {n.label}
                </Link>
              );
            })}
          </nav>
          <button
            onClick={logout}
            className="flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-muted"
          >
            <LogOut className="size-4" />
            <span className="hidden sm:inline">Logout</span>
          </button>
        </div>
        <nav className="md:hidden border-t border-border overflow-x-auto">
          <div className="container flex gap-1 py-2">
            {nav.map((n) => {
              const active = pathname === n.to;
              const Icon = n.icon;
              return (
                <Link
                  key={n.to}
                  to={n.to}
                  className={cn(
                    "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium whitespace-nowrap",
                    active
                      ? "bg-accent text-accent-foreground"
                      : "text-muted-foreground hover:text-foreground hover:bg-muted"
                  )}
                >
                  <Icon className="size-3.5" />
                  {n.label}
                </Link>
              );
            })}
          </div>
        </nav>
      </header>
      <main className="container py-8 mx-auto">{children}</main>
    </div>
  );
}
