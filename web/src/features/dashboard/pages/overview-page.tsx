import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useDashboardStats } from "../hooks/use-dashboard-stats";

const statCards = [
  { key: "totalWebhooksSent", label: "Webhooks Sent" },
  { key: "totalWebhooksFailed", label: "Webhooks Failed" },
  { key: "totalWebhooksQueued", label: "Webhooks Queued" },
  { key: "successRate", label: "Success Rate" },
] as const;

function formatValue(key: string, value: number): string {
  if (key === "successRate") return `${value}%`;
  return value.toLocaleString();
}

export function OverviewPage() {
  const { stats, isLoading, error } = useDashboardStats();

  return (
    <div className="p-8">
      <h2 className="text-lg font-medium mb-6">Overview</h2>

      {error && (
        <div className="bg-destructive/10 text-destructive rounded-lg p-4 mb-6 text-sm">
          {error}
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((card) => (
          <Card key={card.key}>
            <CardHeader>
              <CardTitle className="text-muted-foreground text-sm font-medium">
                {card.label}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className="h-8 w-24 animate-pulse rounded bg-muted" />
              ) : (
                <p className="text-2xl font-bold">
                  {stats ? formatValue(card.key, stats[card.key]) : "—"}
                </p>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
