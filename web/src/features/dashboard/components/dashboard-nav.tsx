import { NavLink } from "react-router";
import { BarChart3, Key, ScrollText } from "lucide-react";
import { cn } from "@/lib/utils";

const navItems = [
  { to: "/dashboard", label: "Overview", icon: BarChart3, end: true },
  { to: "/dashboard/api-keys", label: "API Keys", icon: Key, end: false },
  { to: "/dashboard/logs", label: "Logs", icon: ScrollText, end: false },
];

export function DashboardNav() {
  return (
    <nav className="flex items-center gap-1 px-8 py-2">
      {navItems.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          end={item.end}
          className={({ isActive }) =>
            cn(
              "inline-flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors",
              isActive
                ? "bg-accent text-accent-foreground"
                : "text-muted-foreground hover:bg-muted hover:text-foreground",
            )
          }
        >
          <item.icon className="size-4" />
          {item.label}
        </NavLink>
      ))}
    </nav>
  );
}
