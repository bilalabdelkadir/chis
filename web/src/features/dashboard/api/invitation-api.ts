import { apiRequest } from "@/shared/api/api-client";
import type {
  Invitation,
  CreateInvitationRequest,
  AcceptInvitationResponse,
  PendingInvitation,
} from "../types/invitation.types";

export function fetchInvitations(): Promise<Invitation[]> {
  return apiRequest<Invitation[]>({
    method: "GET",
    path: "/api/invitations/",
    auth: true,
  });
}

export function createInvitation(
  data: CreateInvitationRequest,
): Promise<Invitation> {
  return apiRequest<Invitation>({
    method: "POST",
    path: "/api/invitations/",
    body: data,
    auth: true,
  });
}

export function cancelInvitation(id: string): Promise<void> {
  return apiRequest<void>({
    method: "DELETE",
    path: `/api/invitations/${id}`,
    auth: true,
  });
}

export function acceptInvitation(
  token: string,
): Promise<AcceptInvitationResponse> {
  return apiRequest<AcceptInvitationResponse>({
    method: "POST",
    path: "/api/invitations/accept",
    body: { token },
    auth: true,
  });
}

export function fetchPendingInvitations(): Promise<PendingInvitation[]> {
  return apiRequest<PendingInvitation[]>({
    method: "GET",
    path: "/api/invitations/pending",
    auth: true,
  });
}

export function respondToInvitation(
  id: string,
  action: "accept" | "decline",
): Promise<AcceptInvitationResponse> {
  return apiRequest<AcceptInvitationResponse>({
    method: "POST",
    path: `/api/invitations/${id}/respond`,
    body: { action },
    auth: true,
  });
}
