import { useCallback, useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  fetchPendingInvitations,
  respondToInvitation,
} from "../api/invitation-api";
import { useOrg } from "@/shared/context/org-context";
import type { PendingInvitation } from "../types/invitation.types";

export function PendingInvitationsDialog() {
  const { addOrg } = useOrg();
  const [invitations, setInvitations] = useState<PendingInvitation[]>([]);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState<Record<string, boolean>>({});

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const data = await fetchPendingInvitations();
        if (!cancelled && data.length > 0) {
          setInvitations(data);
          setOpen(true);
        }
      } catch {
        // silently fail
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleRespond = useCallback(
    async (id: string, action: "accept" | "decline") => {
      setLoading((prev) => ({ ...prev, [id]: true }));
      try {
        const data = await respondToInvitation(id, action);
        if (action === "accept") {
          addOrg({
            id: data.orgId,
            name: data.orgName,
            slug: data.orgSlug,
            role: data.role,
          });
        }
        setInvitations((prev) => {
          const next = prev.filter((inv) => inv.id !== id);
          if (next.length === 0) {
            setOpen(false);
          }
          return next;
        });
      } catch {
        // silently fail
      } finally {
        setLoading((prev) => ({ ...prev, [id]: false }));
      }
    },
    [addOrg],
  );

  if (invitations.length === 0) return null;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Pending invitations</DialogTitle>
          <DialogDescription>
            You have been invited to join the following organizations.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-3">
          {invitations.map((inv) => (
            <div
              key={inv.id}
              className="flex items-center justify-between gap-3 rounded-lg border p-3"
            >
              <div className="min-w-0">
                <p className="truncate text-sm font-medium">{inv.orgName}</p>
                <p className="text-muted-foreground text-xs">
                  Role: {inv.role}
                </p>
              </div>
              <div className="flex shrink-0 gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  disabled={loading[inv.id]}
                  onClick={() => handleRespond(inv.id, "decline")}
                >
                  Decline
                </Button>
                <Button
                  size="sm"
                  disabled={loading[inv.id]}
                  onClick={() => handleRespond(inv.id, "accept")}
                >
                  Accept
                </Button>
              </div>
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}
