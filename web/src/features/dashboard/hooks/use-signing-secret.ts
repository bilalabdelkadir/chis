import { useEffect, useState } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import { useOrg } from "@/shared/context/org-context";
import {
  getSigningSecret,
  rotateSigningSecret,
} from "../api/organization-api";

export function useSigningSecret() {
  const { currentOrg } = useOrg();
  const [signingSecret, setSigningSecret] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRotating, setIsRotating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!currentOrg) return;

    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await getSigningSecret();
        if (!cancelled) setSigningSecret(data.signingSecret);
      } catch (err) {
        if (!cancelled) {
          if (err instanceof ApiRequestError) {
            setError(err.message);
          } else {
            setError("An unexpected error occurred");
          }
        }
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, [currentOrg?.id]);

  async function rotate() {
    setIsRotating(true);
    setError(null);

    try {
      const data = await rotateSigningSecret();
      setSigningSecret(data.signingSecret);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsRotating(false);
    }
  }

  return {
    signingSecret,
    isLoading,
    error,
    rotate,
    isRotating,
  };
}
