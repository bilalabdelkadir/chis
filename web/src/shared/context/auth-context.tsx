import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";
import type { AuthUser } from "@/features/auth/types/auth.types";
import { removeStoredToken, setStoredToken } from "@/shared/api/api-client";

const AUTH_USER_KEY = "auth_user";

interface AuthContextValue {
  token: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  setAuth: (token: string, user: AuthUser) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function getInitialToken(): string | null {
  return localStorage.getItem("auth_token");
}

function getInitialUser(): AuthUser | null {
  const raw = localStorage.getItem(AUTH_USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(getInitialToken);
  const [user, setUser] = useState<AuthUser | null>(getInitialUser);

  const isAuthenticated = token !== null && user !== null;

  const setAuth = useCallback((newToken: string, newUser: AuthUser) => {
    setStoredToken(newToken);
    localStorage.setItem(AUTH_USER_KEY, JSON.stringify(newUser));
    setToken(newToken);
    setUser(newUser);
  }, []);

  const logout = useCallback(() => {
    removeStoredToken();
    localStorage.removeItem(AUTH_USER_KEY);
    localStorage.removeItem("selected_org");
    setToken(null);
    setUser(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ token, user, isAuthenticated, setAuth, logout }),
    [token, user, isAuthenticated, setAuth, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
