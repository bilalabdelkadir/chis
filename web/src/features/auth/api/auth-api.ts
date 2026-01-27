import { apiRequest } from "@/shared/api/api-client";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from "../types/auth.types";

export function loginUser(data: LoginRequest): Promise<LoginResponse> {
  return apiRequest<LoginResponse>({
    method: "POST",
    path: "/auth/login",
    body: data,
    auth: false,
  });
}

export function registerUser(data: RegisterRequest): Promise<RegisterResponse> {
  return apiRequest<RegisterResponse>({
    method: "POST",
    path: "/auth/register",
    body: data,
    auth: false,
  });
}
