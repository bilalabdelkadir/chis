import { useEffect, useState, useRef } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import { useOrg } from "@/shared/context/org-context";
import { fetchWebhookLogs } from "../api/dashboard-api";
import type { WebhookLog } from "../types/dashboard.types";

export function useWebhookLogs() {
  const { currentOrg } = useOrg();
  const [logs, setLogs] = useState<WebhookLog[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  // Reset page when filters change
  const isFirstRender = useRef(true);
  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }
    setPage(1);
  }, [statusFilter, debouncedSearch]);

  // Fetch logs
  useEffect(() => {
    if (!currentOrg) return;

    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const response = await fetchWebhookLogs({
          page,
          limit: 20,
          status: statusFilter || undefined,
          search: debouncedSearch || undefined,
        });
        if (!cancelled) {
          setLogs(response.data);
          setTotalPages(response.totalPages);
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
  }, [page, statusFilter, debouncedSearch, currentOrg?.id]);

  function goToPage(p: number) {
    setPage(p);
  }

  return {
    logs,
    isLoading,
    error,
    page,
    totalPages,
    goToPage,
    statusFilter,
    setStatusFilter,
    search,
    setSearch,
  };
}
