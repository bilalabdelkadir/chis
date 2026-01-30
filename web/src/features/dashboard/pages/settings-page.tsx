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
import { useOrg } from "@/shared/context/org-context";
import { deleteOrganization } from "../api/organization-api";

export function SettingsPage() {
  const { currentOrg, removeOrg } = useOrg();
  const navigate = useNavigate();
  const [confirmName, setConfirmName] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);

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

  return (
    <div className="p-8">
      <div className="mb-6">
        <h2 className="text-lg font-medium">Settings</h2>
        <p className="text-muted-foreground text-sm mt-1">
          Manage your organization settings.
        </p>
      </div>

      <div className="space-y-6">
        <div>
          <h3 className="text-sm font-medium text-destructive mb-3">
            Danger zone
          </h3>
          <Card className="border-destructive/30">
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
        </div>
      </div>
    </div>
  );
}
