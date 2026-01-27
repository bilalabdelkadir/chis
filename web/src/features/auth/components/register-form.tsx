import { Link } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { useRegister } from "../hooks/use-register";

export function RegisterForm() {
  const {
    firstName,
    setFirstName,
    lastName,
    setLastName,
    email,
    setEmail,
    password,
    setPassword,
    error,
    isLoading,
    getFieldError,
    handleSubmit,
  } = useRegister();

  return (
    <form onSubmit={handleSubmit}>
      <FieldGroup className="gap-6">
        {error && (
          <div className="bg-destructive/10 text-destructive rounded-lg px-4 py-3 text-sm font-medium">{error}</div>
        )}
        <div className="grid grid-cols-2 gap-4">
          <Field>
            <FieldLabel>First name</FieldLabel>
            <Input
              type="text"
              placeholder="John"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              required
            />
            {getFieldError("firstName") && (
              <FieldError>{getFieldError("firstName")}</FieldError>
            )}
          </Field>
          <Field>
            <FieldLabel>Last name</FieldLabel>
            <Input
              type="text"
              placeholder="Doe"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              required
            />
            {getFieldError("lastName") && (
              <FieldError>{getFieldError("lastName")}</FieldError>
            )}
          </Field>
        </div>
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
            placeholder="Create a password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {getFieldError("password") && (
            <FieldError>{getFieldError("password")}</FieldError>
          )}
        </Field>
        <Button type="submit" disabled={isLoading} className="w-full">
          {isLoading ? "Creating account..." : "Create account"}
        </Button>
        <p className="text-muted-foreground text-center text-sm">
          Already have an account?{" "}
          <Link to="/login" className="text-primary underline underline-offset-4 hover:text-primary/80">
            Sign in
          </Link>
        </p>
      </FieldGroup>
    </form>
  );
}
