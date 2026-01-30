import { ApiRequestError, type ApiErrorBody } from './api-error';

const AUTH_TOKEN_KEY = 'auth_token';
const SELECTED_ORG_KEY = 'selected_org';
const BASE_URL = import.meta.env.VITE_API_URL ?? '';

export function getStoredToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function setStoredToken(token: string): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function removeStoredToken(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}

interface ApiRequestOptions {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  path: string;
  body?: unknown;
  auth?: boolean;
}

export async function apiRequest<T>({
  method,
  path,
  body,
  auth = false,
}: ApiRequestOptions): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (auth) {
    const token = getStoredToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    const orgId = localStorage.getItem(SELECTED_ORG_KEY);
    if (orgId) {
      headers['X-Org-ID'] = orgId;
    }
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    const errorBody: ApiErrorBody = await response.json();
    throw new ApiRequestError(response.status, errorBody);
  }

  if (response.status === 204) return undefined as T;

  return response.json() as Promise<T>;
}
