import { Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useWebhookLogs } from "../hooks/use-webhook-logs";
import type { WebhookLog } from "../types/dashboard.types";

const statusBadgeVariant: Record<
  WebhookLog["status"],
  "default" | "destructive" | "secondary"
> = {
  success: "default",
  failed: "destructive",
  pending: "secondary",
};

export function LogsPage() {
  const {
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
  } = useWebhookLogs();

  return (
    <div className="p-8">
      <div className="mb-6">
        <h2 className="text-lg font-medium">Webhook Logs</h2>
        <p className="text-muted-foreground text-sm mt-1">
          View delivery logs for your webhooks.
        </p>
      </div>

      {error && (
        <div className="bg-destructive/10 text-destructive rounded-lg p-4 mb-6 text-sm">
          {error}
        </div>
      )}

      <div className="flex items-center gap-3 mb-6">
        <div className="relative max-w-xs flex-1">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 size-4 text-muted-foreground pointer-events-none" />
          <Input
            placeholder="Search endpoints..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>
        <Select
          value={statusFilter}
          onValueChange={(val) => setStatusFilter(val === "all" ? "" : (val ?? ""))}
        >
          <SelectTrigger className="w-36">
            <SelectValue placeholder="All statuses" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All statuses</SelectItem>
            <SelectItem value="success">Success</SelectItem>
            <SelectItem value="failed">Failed</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Card>
        <CardContent className="p-0">
          {/* Header row */}
          <div className="grid grid-cols-6 gap-4 px-4 py-2 text-xs text-muted-foreground font-medium border-b border-border">
            <span>Status</span>
            <span>Endpoint</span>
            <span>Event</span>
            <span>Code</span>
            <span>Time</span>
            <span>Date</span>
          </div>

          {isLoading ? (
            <div className="divide-y divide-border">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="px-4 py-3">
                  <div className="h-5 animate-pulse rounded bg-muted" />
                </div>
              ))}
            </div>
          ) : logs.length === 0 ? (
            <div className="px-4 py-8 text-center text-muted-foreground text-sm">
              No logs found.
            </div>
          ) : (
            <div className="divide-y divide-border">
              {logs.map((log) => (
                <div
                  key={log.id}
                  className="grid grid-cols-6 gap-4 px-4 py-3 text-sm items-center"
                >
                  <span>
                    <Badge variant={statusBadgeVariant[log.status]}>
                      {log.status}
                    </Badge>
                  </span>
                  <span className="font-mono text-xs truncate">
                    {log.endpoint}
                  </span>
                  <span className="text-muted-foreground text-xs">
                    {log.eventType}
                  </span>
                  <span className="font-mono text-xs">{log.statusCode}</span>
                  <span className="text-muted-foreground text-xs">
                    {log.responseTimeMs}ms
                  </span>
                  <span className="text-muted-foreground text-xs">
                    {new Date(log.attemptedAt).toLocaleDateString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm text-muted-foreground">
            Page {page} of {totalPages}
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() => goToPage(page - 1)}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= totalPages}
              onClick={() => goToPage(page + 1)}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
