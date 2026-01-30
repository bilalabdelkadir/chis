import { apiRequest } from "@/shared/api/api-client";
import type {
  DashboardStats,
  ApiKey,
  CreateApiKeyRequest,
  CreateApiKeyResponse,
  WebhookLog,
  WebhookLogDetail,
  WebhookLogsParams,
  PaginatedResponse,
  SendWebhookRequest,
  SendWebhookResponse,
} from "../types/dashboard.types";

export function fetchDashboardStats(): Promise<DashboardStats> {
  return apiRequest<DashboardStats>({
    method: "GET",
    path: "/api/dashboard/stats",
    auth: true,
  });
}

export function fetchApiKeys(): Promise<ApiKey[]> {
  return apiRequest<ApiKey[]>({
    method: "GET",
    path: "/api/api-key/list",
    auth: true,
  });
}

export function createApiKey(
  data: CreateApiKeyRequest,
): Promise<CreateApiKeyResponse> {
  return apiRequest<CreateApiKeyResponse>({
    method: "POST",
    path: "/api/api-key/create",
    body: data,
    auth: true,
  });
}

export function deleteApiKey(id: string): Promise<void> {
  return apiRequest<void>({
    method: "DELETE",
    path: `/api/api-key/${id}`,
    auth: true,
  });
}

export function fetchWebhookLogDetail(id: string): Promise<WebhookLogDetail> {
  return apiRequest<WebhookLogDetail>({
    method: "GET",
    path: `/api/webhook-logs/${id}`,
    auth: true,
  });
}

export function fetchWebhookLogs(
  params?: WebhookLogsParams,
): Promise<PaginatedResponse<WebhookLog>> {
  const searchParams = new URLSearchParams();
  if (params?.status) searchParams.set("status", params.status);
  if (params?.search) searchParams.set("search", params.search);
  if (params?.page) searchParams.set("page", String(params.page));
  if (params?.limit) searchParams.set("limit", String(params.limit));

  const query = searchParams.toString();
  const path = query ? `/api/webhook-logs?${query}` : "/api/webhook-logs";

  return apiRequest<PaginatedResponse<WebhookLog>>({
    method: "GET",
    path,
    auth: true,
  });
}

const BASE_URL = import.meta.env.VITE_API_URL ?? "";

export async function sendTestWebhook(
  apiKey: string,
  request: SendWebhookRequest,
): Promise<SendWebhookResponse> {
  const response = await fetch(`${BASE_URL}/webhook/send`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": apiKey,
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error || `Request failed with status ${response.status}`);
  }

  return response.json() as Promise<SendWebhookResponse>;
}
