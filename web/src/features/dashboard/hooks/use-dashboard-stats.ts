import { useEffect, useState } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import { useOrg } from "@/shared/context/org-context";
import { fetchDashboardStats } from "../api/dashboard-api";
import type { DashboardStats } from "../types/dashboard.types";

export function useDashboardStats() {
  const { currentOrg } = useOrg();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!currentOrg) return;

    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await fetchDashboardStats();
        if (!cancelled) setStats(data);
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

  return { stats, isLoading, error };
}
