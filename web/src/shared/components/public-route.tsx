import { Navigate, Outlet } from "react-router";
import { useAuth } from "@/shared/hooks/use-auth";

export function PublicRoute() {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
