import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card";
import { PageHeader } from "../components/page-header";
import { useDashboardStats } from "../hooks/use-dashboard-stats";

const statCards = [
  { key: "totalWebhooksSent", label: "Webhooks sent", description: "Total deliveries" },
  { key: "totalWebhooksFailed", label: "Failed", description: "Requires attention" },
  { key: "totalWebhooksQueued", label: "Queued", description: "Pending delivery" },
  { key: "successRate", label: "Success rate", description: "Last 30 days" },
] as const;

function formatValue(key: string, value: number): string {
  if (key === "successRate") return `${value}%`;
  return value.toLocaleString();
}

export function OverviewPage() {
  const { stats, isLoading, error } = useDashboardStats();

  return (
    <>
      <PageHeader
        title="Overview"
        description="Key metrics for your webhook delivery and API usage."
      />

      {error && (
        <div className="mb-6 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((card) => (
          <Card key={card.key} className="border-border/80 transition-colors hover:border-border">
            <CardHeader className="pb-2">
              <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
                {card.label}
              </p>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="h-9 w-24 animate-pulse rounded-md bg-muted" />
              ) : (
                <>
                  <p className="text-2xl font-semibold tabular-nums tracking-tight text-foreground">
                    {stats ? formatValue(card.key, stats[card.key]) : "—"}
                  </p>
                  <p className="mt-0.5 text-xs text-muted-foreground">
                    {card.description}
                  </p>
                </>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
}
