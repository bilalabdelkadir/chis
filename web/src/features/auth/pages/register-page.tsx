import { AuthLayout } from "../components/auth-layout";
import { RegisterForm } from "../components/register-form";

export function RegisterPage() {
  return (
    <AuthLayout
      title="Create an account"
      description="Enter your details to create a new account"
    >
      <RegisterForm />
    </AuthLayout>
  );
}
