import { useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { PageHeader } from "../components/page-header";
import { useOrg } from "@/shared/context/org-context";
import { deleteOrganization } from "../api/organization-api";
import { useSigningSecret } from "../hooks/use-signing-secret";

export function SettingsPage() {
  const { currentOrg, removeOrg } = useOrg();
  const navigate = useNavigate();
  const [confirmName, setConfirmName] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);

  const {
    signingSecret,
    isLoading: isSecretLoading,
    error: secretError,
    rotate,
    isRotating,
  } = useSigningSecret();

  const [secretVisible, setSecretVisible] = useState(false);
  const [copied, setCopied] = useState(false);
  const [rotateOpen, setRotateOpen] = useState(false);

  const orgName = currentOrg?.name ?? "";
  const canDelete = confirmName === orgName;

  async function handleDelete() {
    if (!currentOrg || !canDelete) return;
    setIsDeleting(true);
    setError("");
    try {
      await deleteOrganization(confirmName);
      const orgId = currentOrg.id;
      removeOrg(orgId);
      setOpen(false);
      navigate("/dashboard", { replace: true });
    } catch {
      setError("Failed to delete organization. Please try again.");
    } finally {
      setIsDeleting(false);
    }
  }

  function handleCopy() {
    if (signingSecret) {
      navigator.clipboard.writeText(signingSecret);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  }

  async function handleRotate() {
    await rotate();
    setRotateOpen(false);
  }

  function maskedSecret(secret: string) {
    if (secret.length <= 6) return secret;
    return secret.slice(0, 6) + "\u2022".repeat(12);
  }

  return (
    <>
      <PageHeader
        title="Settings"
        description="Organization settings, signing secret, and danger zone."
      />

      <div className="space-y-8">
        <section>
          <h2 className="text-sm font-medium text-foreground mb-3">Webhook Signing Secret</h2>
          <Card className="border-border/80">
            <CardContent className="py-4">
              <p className="text-muted-foreground text-xs mb-3">
                Used to sign webhook payloads so receivers can verify
                authenticity. Every delivery includes{" "}
                <code className="text-xs bg-muted px-1 py-0.5 rounded">
                  X-Webhook-Signature
                </code>{" "}
                headers.
              </p>

              {isSecretLoading ? (
                <div className="h-9 bg-muted animate-pulse rounded" />
              ) : secretError ? (
                <p className="text-destructive text-xs">{secretError}</p>
              ) : signingSecret ? (
                <div className="flex items-center gap-2">
                  <code className="flex-1 text-sm bg-muted px-3 py-2 rounded font-mono select-all">
                    {secretVisible
                      ? signingSecret
                      : maskedSecret(signingSecret)}
                  </code>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSecretVisible(!secretVisible)}
                  >
                    {secretVisible ? "Hide" : "Reveal"}
                  </Button>
                  <Button variant="outline" size="sm" onClick={handleCopy}>
                    {copied ? "Copied!" : "Copy"}
                  </Button>
                  <AlertDialog
                    open={rotateOpen}
                    onOpenChange={setRotateOpen}
                  >
                    <AlertDialogTrigger
                      render={
                        <Button variant="outline" size="sm">
                          Rotate
                        </Button>
                      }
                    />
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          Rotate signing secret?
                        </AlertDialogTitle>
                        <AlertDialogDescription>
                          The current secret will be immediately invalidated.
                          All webhook receivers using the old secret will fail
                          verification until updated.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <Button
                          variant="destructive"
                          disabled={isRotating}
                          onClick={handleRotate}
                        >
                          {isRotating ? "Rotating..." : "Rotate secret"}
                        </Button>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              ) : null}
            </CardContent>
          </Card>
        </section>

        <section>
          <h2 className="text-sm font-medium text-destructive mb-3">
            Danger zone
          </h2>
          <Card className="border-destructive/30 border-border/80">
            <CardContent className="flex items-center justify-between py-4">
              <div>
                <p className="text-sm font-medium">Delete organization</p>
                <p className="text-muted-foreground text-xs mt-0.5">
                  Permanently delete this organization and all its data.
                </p>
              </div>
              <AlertDialog
                open={open}
                onOpenChange={(isOpen) => {
                  setOpen(isOpen);
                  if (!isOpen) {
                    setConfirmName("");
                    setError("");
                  }
                }}
              >
                <AlertDialogTrigger
                  render={
                    <Button variant="destructive" size="sm">
                      Delete organization
                    </Button>
                  }
                />
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete {orgName}?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This action cannot be undone. All API keys, webhooks, and
                      members will be permanently removed.
                    </AlertDialogDescription>
                  </AlertDialogHeader>

                  <div className="px-0.5">
                    <label className="text-sm text-muted-foreground">
                      Type <span className="font-semibold text-foreground">{orgName}</span> to
                      confirm
                    </label>
                    <Input
                      className="mt-1.5"
                      value={confirmName}
                      onChange={(e) => setConfirmName(e.target.value)}
                      placeholder={orgName}
                    />
                    {error && (
                      <p className="text-destructive text-xs mt-2">{error}</p>
                    )}
                  </div>

                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <Button
                      variant="destructive"
                      disabled={!canDelete || isDeleting}
                      onClick={handleDelete}
                    >
                      {isDeleting ? "Deleting..." : "Delete"}
                    </Button>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </CardContent>
          </Card>
        </section>
      </div>
    </>
  );
}
