export interface Organization {
  id: string;
  name: string;
  slug: string;
  role: string;
}

export interface CreateOrganizationRequest {
  name: string;
}
