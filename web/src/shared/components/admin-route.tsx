import { Navigate, Outlet } from "react-router";
import { useOrg } from "@/shared/context/org-context";

export function AdminRoute() {
  const { currentOrg } = useOrg();

  if (currentOrg?.role !== "admin") {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
