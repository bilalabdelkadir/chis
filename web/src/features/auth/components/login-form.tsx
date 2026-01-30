import { Link, useSearchParams } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { useLogin } from "../hooks/use-login";

export function LoginForm() {
  const [searchParams] = useSearchParams();
  const redirect = searchParams.get("redirect");

  const {
    email,
    setEmail,
    password,
    setPassword,
    error,
    isLoading,
    getFieldError,
    handleSubmit,
  } = useLogin();

  return (
    <form onSubmit={handleSubmit}>
      <FieldGroup className="gap-6">
        {error && (
          <div className="bg-destructive/10 text-destructive rounded-lg px-4 py-3 text-sm font-medium">{error}</div>
        )}
        <Field>
          <FieldLabel>Email</FieldLabel>
          <Input
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          {getFieldError("email") && (
            <FieldError>{getFieldError("email")}</FieldError>
          )}
        </Field>
        <Field>
          <FieldLabel>Password</FieldLabel>
          <Input
            type="password"
            placeholder="Enter your password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {getFieldError("password") && (
            <FieldError>{getFieldError("password")}</FieldError>
          )}
        </Field>
        <Button type="submit" disabled={isLoading} className="w-full">
          {isLoading ? "Signing in..." : "Sign in"}
        </Button>
        <p className="text-muted-foreground text-center text-sm">
          Don't have an account?{" "}
          <Link to={redirect ? `/register?redirect=${encodeURIComponent(redirect)}` : "/register"} className="text-primary underline underline-offset-4 hover:text-primary/80">
            Register
          </Link>
        </p>
      </FieldGroup>
    </form>
  );
}
