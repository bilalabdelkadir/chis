export interface Invitation {
  id: string;
  email: string;
  role: string;
  status: string;
  expiresAt: string;
  createdAt: string;
}

export interface CreateInvitationRequest {
  email: string;
}

export interface AcceptInvitationResponse {
  message: string;
  orgId: string;
  orgName: string;
  orgSlug: string;
  role: string;
}

export interface PendingInvitation {
  id: string;
  orgName: string;
  orgSlug: string;
  role: string;
  expiresAt: string;
}
