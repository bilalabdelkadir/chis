import { useState } from "react";
import { Send, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { PageHeader } from "../components/page-header";
import { useInvitations } from "../hooks/use-invitations";

function statusVariant(status: string) {
  switch (status) {
    case "accepted":
      return "secondary";
    case "expired":
    case "cancelled":
      return "outline";
    default:
      return "default";
  }
}

export function MembersPage() {
  const { invitations, isLoading, isSending, error, sendInvitation, cancelInv } =
    useInvitations();

  const [email, setEmail] = useState("");

  async function handleSend(e: React.FormEvent) {
    e.preventDefault();
    if (!email.trim()) return;
    await sendInvitation(email.trim());
    setEmail("");
  }

  return (
    <>
      <PageHeader
        title="Members"
        description="Invite people to your organization and manage pending invitations."
      />

      {error && (
        <div className="mb-6 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <form onSubmit={handleSend} className="flex items-center gap-3 mb-6">
        <Input
          type="email"
          placeholder="Email address"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="max-w-xs"
        />
        <Button type="submit" disabled={isSending || !email.trim()}>
          <Send className="size-4" />
          Send invite
        </Button>
      </form>

      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-16 animate-pulse rounded-xl bg-muted"
            />
          ))}
        </div>
      ) : invitations.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <p className="text-sm font-medium text-muted-foreground">No invitations yet</p>
            <p className="mt-1 text-xs text-muted-foreground">Invite teammates by email above.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {invitations.map((inv) => (
            <Card key={inv.id} className="border-border/80">
              <CardContent className="flex items-center justify-between py-4">
                <div className="flex items-center gap-4">
                  <div>
                    <p className="text-sm font-medium">{inv.email}</p>
                    <p className="text-muted-foreground text-xs">
                      Invited {new Date(inv.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <Badge variant={statusVariant(inv.status)}>{inv.status}</Badge>
                </div>
                {inv.status === "pending" && (
                  <AlertDialog>
                    <AlertDialogTrigger
                      render={
                        <Button variant="ghost" size="icon-sm">
                          <X className="size-4 text-muted-foreground" />
                        </Button>
                      }
                    />
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Cancel invitation</AlertDialogTitle>
                        <AlertDialogDescription>
                          Are you sure you want to cancel the invitation to{" "}
                          {inv.email}? They will no longer be able to join using
                          this invitation.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Keep</AlertDialogCancel>
                        <AlertDialogAction
                          variant="destructive"
                          onClick={() => cancelInv(inv.id)}
                        >
                          Cancel invitation
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </>
  );
}
