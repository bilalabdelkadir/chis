import { useState } from "react";
import { Plus, Trash2, Copy, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
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
import { useApiKeys } from "../hooks/use-api-keys";

export function ApiKeysPage() {
  const {
    apiKeys,
    isLoading,
    isCreating,
    error,
    createKey,
    removeKey,
    newlyCreatedKey,
    clearNewKey,
  } = useApiKeys();

  const [newKeyName, setNewKeyName] = useState("");

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!newKeyName.trim()) return;
    await createKey(newKeyName.trim());
    setNewKeyName("");
  }

  function handleCopy() {
    if (newlyCreatedKey) {
      navigator.clipboard.writeText(newlyCreatedKey.key);
    }
  }

  return (
    <>
      <PageHeader
        title="API Keys"
        description="Create and manage API keys for programmatic access to the webhook API."
      />

      {error && (
        <div className="mb-6 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {error}
        </div>
      )}

      {newlyCreatedKey && (
        <Card className="mb-6 border-primary/30 bg-primary/5">
          <CardContent className="pt-4">
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium mb-2">
                  Your new API key has been created. Copy it now — you won't be
                  able to see it again.
                </p>
                <code className="bg-muted text-sm px-3 py-2 rounded-lg block break-all">
                  {newlyCreatedKey.key}
                </code>
              </div>
              <div className="flex items-center gap-1 shrink-0">
                <Button variant="ghost" size="icon-sm" onClick={handleCopy}>
                  <Copy className="size-4" />
                </Button>
                <Button variant="ghost" size="icon-sm" onClick={clearNewKey}>
                  <X className="size-4" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <form onSubmit={handleCreate} className="flex items-center gap-3 mb-6">
        <Input
          placeholder="Key name"
          value={newKeyName}
          onChange={(e) => setNewKeyName(e.target.value)}
          className="max-w-xs"
        />
        <Button type="submit" disabled={isCreating || !newKeyName.trim()}>
          <Plus className="size-4" />
          Create key
        </Button>
      </form>

      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-16 animate-pulse rounded-lg bg-muted"
            />
          ))}
        </div>
      ) : apiKeys.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <p className="text-sm font-medium text-muted-foreground">No API keys yet</p>
            <p className="mt-1 text-xs text-muted-foreground">Create a key to start sending webhooks.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {apiKeys.map((key) => (
            <Card key={key.id} className="border-border/80">
              <CardContent className="flex items-center justify-between py-4">
                <div className="flex items-center gap-6">
                  <div>
                    <p className="text-sm font-medium">{key.name}</p>
                    <code className="text-muted-foreground text-xs">
                      {key.prefix}...
                    </code>
                  </div>
                  <p className="text-muted-foreground text-xs">
                    {key.lastUsedAt
                      ? `Last used ${new Date(key.lastUsedAt).toLocaleDateString()}`
                      : "Never used"}
                  </p>
                </div>
                <AlertDialog>
                  <AlertDialogTrigger
                    render={
                      <Button variant="ghost" size="icon-sm">
                        <Trash2 className="size-4 text-muted-foreground" />
                      </Button>
                    }
                  />
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Revoke API key</AlertDialogTitle>
                      <AlertDialogDescription>
                        Are you sure you want to revoke "{key.name}"? This
                        action cannot be undone. Any integrations using this key
                        will stop working.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        variant="destructive"
                        onClick={() => removeKey(key.id)}
                      >
                        Revoke
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </>
  );
}
