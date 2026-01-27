import { useEffect, useState } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import {
  fetchApiKeys,
  createApiKey,
  deleteApiKey,
} from "../api/dashboard-api";
import type { ApiKey, CreateApiKeyResponse } from "../types/dashboard.types";

export function useApiKeys() {
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [newlyCreatedKey, setNewlyCreatedKey] =
    useState<CreateApiKeyResponse | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await fetchApiKeys();
        if (!cancelled) setApiKeys(data);
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
  }, []);

  async function createKey(name: string) {
    setIsCreating(true);
    setError(null);

    try {
      const response = await createApiKey({ name });
      setNewlyCreatedKey(response);
      setApiKeys((prev) => [
        {
          id: response.id,
          name: response.name,
          prefix: response.prefix,
          createdAt: response.createdAt,
          lastUsedAt: null,
        },
        ...prev,
      ]);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsCreating(false);
    }
  }

  async function removeKey(id: string) {
    setError(null);

    try {
      await deleteApiKey(id);
      setApiKeys((prev) => prev.filter((k) => k.id !== id));
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    }
  }

  function clearNewKey() {
    setNewlyCreatedKey(null);
  }

  return {
    apiKeys,
    isLoading,
    isCreating,
    error,
    createKey,
    removeKey,
    newlyCreatedKey,
    clearNewKey,
  };
}
