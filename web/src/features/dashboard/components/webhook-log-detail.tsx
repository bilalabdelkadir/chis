import { Badge } from '@/components/ui/badge';
import { useWebhookLogDetail } from '../hooks/use-webhook-log-detail';
import type { DeliveryAttemptDetail } from '../types/dashboard.types';

const statusBadgeVariant: Record<
  string,
  'default' | 'destructive' | 'secondary'
> = {
  success: 'default',
  failed: 'destructive',
  pending: 'secondary',
};

function formatJson(value: unknown): string {
  if (value == null) return '';
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }
  return JSON.stringify(value, null, 2);
}

function StatusCodeBadge({ code }: { code: number | null }) {
  if (code == null)
    return <span className="text-muted-foreground text-xs">-</span>;
  const variant = code >= 200 && code < 300 ? 'default' : 'destructive';
  return <Badge variant={variant}>{code}</Badge>;
}

function DetailSkeleton() {
  return (
    <div className="space-y-6 p-6">
      <div className="space-y-2">
        <div className="h-6 w-48 animate-pulse rounded bg-muted" />
        <div className="h-4 w-full animate-pulse rounded bg-muted" />
      </div>
      <div className="grid grid-cols-2 gap-4">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="space-y-1">
            <div className="h-3 w-16 animate-pulse rounded bg-muted" />
            <div className="h-5 w-24 animate-pulse rounded bg-muted" />
          </div>
        ))}
      </div>
      <div className="space-y-2">
        <div className="h-4 w-24 animate-pulse rounded bg-muted" />
        <div className="h-32 w-full animate-pulse rounded bg-muted" />
      </div>
      <div className="space-y-2">
        <div className="h-4 w-32 animate-pulse rounded bg-muted" />
        <div className="h-20 w-full animate-pulse rounded bg-muted" />
      </div>
    </div>
  );
}

function AttemptCard({ attempt }: { attempt: DeliveryAttemptDetail }) {
  const responseFormatted = attempt.responseBody
    ? formatJson(attempt.responseBody)
    : null;

  return (
    <div className="rounded-lg border p-4 space-y-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">
          Attempt #{attempt.attemptNumber}
        </span>
        <StatusCodeBadge code={attempt.statusCode} />
      </div>
      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        {attempt.durationMs != null && <span>{attempt.durationMs}ms</span>}
        <span>{new Date(attempt.attemptedAt).toLocaleString()}</span>
      </div>
      {attempt.errorMessage && (
        <div className="rounded bg-destructive/10 px-3 py-2 text-sm text-destructive">
          {attempt.errorMessage}
        </div>
      )}
      {responseFormatted && (
        <div className="space-y-1">
          <span className="text-xs font-medium text-muted-foreground">
            Response Body
          </span>
          <pre className="overflow-auto rounded bg-muted p-3 text-xs max-h-80">
            {responseFormatted}
          </pre>
        </div>
      )}
    </div>
  );
}

interface WebhookLogDetailViewProps {
  logId: string | null;
}

export function WebhookLogDetailView({ logId }: WebhookLogDetailViewProps) {
  const { detail, isLoading, error } = useWebhookLogDetail(logId);

  if (isLoading) return <DetailSkeleton />;

  if (error) {
    return (
      <div className="p-6">
        <div className="rounded-lg bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      </div>
    );
  }

  if (!detail) return null;

  const payloadFormatted = formatJson(detail.payload);

  return (
    <div className="flex-1 overflow-y-auto p-6 space-y-6">
      {/* Status + Method + URL */}
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <Badge variant={statusBadgeVariant[detail.status] ?? 'secondary'}>
            {detail.status}
          </Badge>
          <span className="rounded bg-muted px-2 py-0.5 text-xs font-mono font-medium">
            {detail.method}
          </span>
        </div>
        <p className="font-mono text-sm break-all">{detail.url}</p>
      </div>

      {/* Metadata grid */}
      <div className="grid grid-cols-2 gap-4 rounded-lg border p-4">
        <div>
          <p className="text-xs text-muted-foreground">Attempts</p>
          <p className="text-sm font-medium">{detail.attemptCount}</p>
        </div>
        <div>
          <p className="text-xs text-muted-foreground">Created</p>
          <p className="text-sm font-medium">
            {new Date(detail.createdAt).toLocaleString()}
          </p>
        </div>
        <div>
          <p className="text-xs text-muted-foreground">Updated</p>
          <p className="text-sm font-medium">
            {new Date(detail.updatedAt).toLocaleString()}
          </p>
        </div>
        <div>
          <p className="text-xs text-muted-foreground">Next Retry</p>
          <p className="text-sm font-medium">
            {detail.nextRetryAt
              ? new Date(detail.nextRetryAt).toLocaleString()
              : '-'}
          </p>
        </div>
      </div>

      {/* Request payload */}
      {payloadFormatted && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Request Payload</h4>
          <pre className="overflow-auto rounded-lg bg-muted p-4 text-xs max-h-96">
            {payloadFormatted}
          </pre>
        </div>
      )}

      {/* Delivery attempts */}
      <div className="space-y-3">
        <h4 className="text-sm font-medium">
          Delivery Attempts ({detail.deliveryAttempts.length})
        </h4>
        {detail.deliveryAttempts.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            No delivery attempts yet.
          </p>
        ) : (
          <div className="space-y-3">
            {detail.deliveryAttempts.map((attempt) => (
              <AttemptCard key={attempt.id} attempt={attempt} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
