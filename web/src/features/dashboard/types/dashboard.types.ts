export interface DashboardStats {
  totalWebhooksSent: number;
  totalWebhooksFailed: number;
  totalWebhooksQueued: number;
  successRate: number;
}

export interface ApiKey {
  id: string;
  name: string;
  prefix: string;
  createdAt: string;
  lastUsedAt: string | null;
}

export interface CreateApiKeyRequest {
  name: string;
}

export interface CreateApiKeyResponse {
  id: string;
  name: string;
  prefix: string;
  key: string;
  createdAt: string;
}

export interface WebhookLog {
  id: string;
  endpoint: string;
  status: "success" | "failed" | "pending";
  statusCode: number;
  eventType: string;
  attemptedAt: string;
  responseTimeMs: number;
}

export interface WebhookLogsParams {
  status?: string;
  search?: string;
  page?: number;
  limit?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}
