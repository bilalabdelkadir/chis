import { Outlet } from "react-router";
import { Separator } from "@/components/ui/separator";
import { DashboardHeader } from "../components/dashboard-header";
import { DashboardNav } from "../components/dashboard-nav";

export function DashboardPage() {
  return (
    <div className="bg-background flex min-h-screen flex-col">
      <DashboardHeader />
      <DashboardNav />
      <Separator />
      <Outlet />
    </div>
  );
}
