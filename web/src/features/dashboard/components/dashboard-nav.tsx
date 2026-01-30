import { NavLink } from "react-router";
import { BarChart3, Key, ScrollText, Play, Users, Settings } from "lucide-react";
import { cn } from "@/lib/utils";
import { useOrg } from "@/shared/context/org-context";

const navItems = [
  { to: "/dashboard", label: "Overview", icon: BarChart3, end: true },
  { to: "/dashboard/api-keys", label: "API Keys", icon: Key, end: false },
  { to: "/dashboard/logs", label: "Logs", icon: ScrollText, end: false },
  { to: "/dashboard/playground", label: "Playground", icon: Play, end: false },
];

const adminNavItems = [
  { to: "/dashboard/members", label: "Members", icon: Users, end: false },
  { to: "/dashboard/settings", label: "Settings", icon: Settings, end: false },
];

export function DashboardNav() {
  const { currentOrg } = useOrg();
  const isAdmin = currentOrg?.role === "admin";

  const items = isAdmin ? [...navItems, ...adminNavItems] : navItems;

  return (
    <nav className="flex items-center gap-1 px-8 py-2">
      {items.map((item) => (
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
