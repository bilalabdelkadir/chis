import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { ApiRequestError, type FieldError } from "@/shared/api/api-error";
import { registerUser } from "../api/auth-api";

export function useRegister() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirect = searchParams.get("redirect");

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<FieldError[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  function getFieldError(field: string): string | undefined {
    return fieldErrors.find((e) => e.field === field)?.message;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setFieldErrors([]);
    setIsLoading(true);

    try {
      await registerUser({ email, password, firstName, lastName });
      navigate(redirect ? `/login?redirect=${encodeURIComponent(redirect)}` : "/login", { replace: true });
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
        if (err.details) {
          setFieldErrors(err.details);
        }
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  }

  return {
    firstName,
    setFirstName,
    lastName,
    setLastName,
    email,
    setEmail,
    password,
    setPassword,
    error,
    fieldErrors,
    isLoading,
    getFieldError,
    handleSubmit,
  };
}
