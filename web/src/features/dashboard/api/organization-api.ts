import { apiRequest } from "@/shared/api/api-client";
import type {
  Organization,
  CreateOrganizationRequest,
} from "../types/organization.types";

export function fetchOrganizations(): Promise<Organization[]> {
  return apiRequest<Organization[]>({
    method: "GET",
    path: "/api/orgs",
    auth: true,
  });
}

export function createOrganization(
  data: CreateOrganizationRequest,
): Promise<Organization> {
  return apiRequest<Organization>({
    method: "POST",
    path: "/api/orgs",
    body: data,
    auth: true,
  });
}
