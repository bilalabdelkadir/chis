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
    <div className="p-8">
      <div className="mb-6">
        <h2 className="text-lg font-medium">API Keys</h2>
        <p className="text-muted-foreground text-sm mt-1">
          Manage your API keys for programmatic access.
        </p>
      </div>

      {error && (
        <div className="bg-destructive/10 text-destructive rounded-lg p-4 mb-6 text-sm">
          {error}
        </div>
      )}

      {newlyCreatedKey && (
        <Card className="border-primary/30 mb-6">
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
              className="h-16 animate-pulse rounded-xl bg-muted"
            />
          ))}
        </div>
      ) : apiKeys.length === 0 ? (
        <Card>
          <CardContent className="py-8">
            <p className="text-muted-foreground text-sm text-center">
              No API keys yet. Create one to get started.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {apiKeys.map((key) => (
            <Card key={key.id}>
              <CardContent className="flex items-center justify-between py-3">
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
    </div>
  );
}
