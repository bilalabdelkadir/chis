import { useEffect, useState } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import { fetchWebhookLogDetail } from "../api/dashboard-api";
import type { WebhookLogDetail } from "../types/dashboard.types";

export function useWebhookLogDetail(logId: string | null) {
  const [detail, setDetail] = useState<WebhookLogDetail | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!logId) {
      setDetail(null);
      setError(null);
      return;
    }

    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await fetchWebhookLogDetail(logId!);
        if (!cancelled) {
          setDetail(data);
        }
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
  }, [logId]);

  return { detail, isLoading, error };
}
