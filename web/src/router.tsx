import { createBrowserRouter, Navigate } from "react-router";
import { PublicRoute } from "@/shared/components/public-route";
import { ProtectedRoute } from "@/shared/components/protected-route";
import { LoginPage } from "@/features/auth/pages/login-page";
import { RegisterPage } from "@/features/auth/pages/register-page";
import { DashboardPage } from "@/features/dashboard/pages/dashboard-page";
import { OverviewPage } from "@/features/dashboard/pages/overview-page";
import { ApiKeysPage } from "@/features/dashboard/pages/api-keys-page";
import { LogsPage } from "@/features/dashboard/pages/logs-page";
import { PlaygroundPage } from "@/features/dashboard/pages/playground-page";
import { MembersPage } from "@/features/dashboard/pages/members-page";
import { AcceptInvitePage } from "@/features/invite/accept-invite-page";

export const router = createBrowserRouter([
  {
    element: <PublicRoute />,
    children: [
      { path: "/login", element: <LoginPage /> },
      { path: "/register", element: <RegisterPage /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: "/dashboard",
        element: <DashboardPage />,
        children: [
          { index: true, element: <OverviewPage /> },
          { path: "api-keys", element: <ApiKeysPage /> },
          { path: "logs", element: <LogsPage /> },
          { path: "playground", element: <PlaygroundPage /> },
          { path: "members", element: <MembersPage /> },
        ],
      },
      { path: "/invite/accept", element: <AcceptInvitePage /> },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/dashboard" replace />,
  },
]);
