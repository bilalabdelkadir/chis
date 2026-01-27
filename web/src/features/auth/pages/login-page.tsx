import { AuthLayout } from "../components/auth-layout";
import { LoginForm } from "../components/login-form";

export function LoginPage() {
  return (
    <AuthLayout
      title="Welcome back"
      description="Enter your credentials to sign in to your account"
    >
      <LoginForm />
    </AuthLayout>
  );
}
