import { RouterProvider } from "react-router";
import { AuthProvider } from "@/shared/context/auth-context";
import { OrgProvider } from "@/shared/context/org-context";
import { router } from "@/router";

export function App() {
  return (
    <AuthProvider>
      <OrgProvider>
        <RouterProvider router={router} />
      </OrgProvider>
    </AuthProvider>
  );
}

export default App;
