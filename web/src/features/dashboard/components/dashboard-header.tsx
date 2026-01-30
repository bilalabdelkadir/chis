import { useNavigate } from "react-router";
import { LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/shared/hooks/use-auth";
import { OrgSwitcher } from "./org-switcher";

export function DashboardHeader() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login", { replace: true });
  }

  return (
    <header>
      <div className="flex items-center justify-between px-8 py-5">
        <div className="flex items-center gap-4">
          <OrgSwitcher />
          <Separator orientation="vertical" className="h-6" />
          <div>
            <p className="text-sm font-medium">
              {user?.firstName || user?.email}
            </p>
            <p className="text-muted-foreground text-xs">{user?.email}</p>
          </div>
        </div>
        <Button variant="ghost" size="sm" onClick={handleLogout}>
          <LogOut className="size-4 text-muted-foreground" />
          Sign out
        </Button>
      </div>
      <Separator />
    </header>
  );
}
