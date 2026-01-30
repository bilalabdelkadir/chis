import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { useAuth } from "@/shared/hooks/use-auth";
import { ApiRequestError, type FieldError } from "@/shared/api/api-error";
import { loginUser } from "../api/auth-api";

export function useLogin() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setAuth } = useAuth();

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
      const response = await loginUser({ email, password });
      setAuth(response.token, { id: "", email, firstName: "" });
      const redirect = searchParams.get("redirect") || "/dashboard";
      navigate(redirect, { replace: true });
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
