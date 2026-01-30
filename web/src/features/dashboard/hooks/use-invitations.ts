import { useEffect, useState } from "react";
import { ApiRequestError } from "@/shared/api/api-error";
import { useOrg } from "@/shared/context/org-context";
import {
  fetchInvitations,
  createInvitation,
  cancelInvitation,
} from "../api/invitation-api";
import type { Invitation } from "../types/invitation.types";

export function useInvitations() {
  const { currentOrg } = useOrg();
  const [invitations, setInvitations] = useState<Invitation[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSending, setIsSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!currentOrg) return;

    let cancelled = false;

    async function load() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await fetchInvitations();
        if (!cancelled) setInvitations(data);
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

  async function sendInvitation(email: string) {
    setIsSending(true);
    setError(null);

    try {
      const invitation = await createInvitation({ email });
      setInvitations((prev) => [invitation, ...prev]);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    } finally {
      setIsSending(false);
    }
  }

  async function cancelInv(id: string) {
    setError(null);

    try {
      await cancelInvitation(id);
      setInvitations((prev) => prev.filter((inv) => inv.id !== id));
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred");
      }
    }
  }

  return {
    invitations,
    isLoading,
    isSending,
    error,
    sendInvitation,
    cancelInv,
  };
}
