import { useState } from "react";
import { Outlet } from "react-router";
import { Menu } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { DashboardSidebar } from "../components/dashboard-sidebar";
import { PendingInvitationsDialog } from "../components/pending-invitations-dialog";

export function DashboardPage() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <div className="flex h-screen min-h-screen flex-col bg-background">
      {/* Mobile: top bar with menu */}
      <header className="flex lg:hidden h-14 shrink-0 items-center gap-2 border-b border-border bg-background px-4">
        <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
          <SheetTrigger>
            <Button variant="ghost" size="icon" className="shrink-0">
              <Menu className="size-5" />
              <span className="sr-only">Open menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-[240px] p-0" showCloseButton>
            <DashboardSidebar onNavigate={() => setMobileOpen(false)} />
          </SheetContent>
        </Sheet>
        <img
          src="/logo.png"
          alt="Chis"
          className="h-7 w-auto dark:invert"
        />
      </header>

      <div className="flex flex-1 min-h-0">
        {/* Desktop: fixed sidebar */}
        <div className="hidden lg:flex lg:flex-col lg:w-[240px] lg:shrink-0 lg:border-r lg:border-border">
          <DashboardSidebar />
        </div>

        {/* Main content */}
        <main className="flex-1 min-w-0 flex flex-col overflow-auto">
          <div className="mx-auto w-full max-w-5xl flex-1 px-4 sm:px-6 lg:px-8 py-6 lg:py-8">
            <Outlet />
          </div>
        </main>
      </div>

      <PendingInvitationsDialog />
    </div>
  );
}
