export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface LoginResponse {
  token: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  firstName: string;
}

export interface AuthUser {
  id: string;
  email: string;
  firstName: string;
}
