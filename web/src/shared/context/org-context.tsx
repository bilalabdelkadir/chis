import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useAuth } from "@/shared/context/auth-context";
import { fetchOrganizations } from "@/features/dashboard/api/organization-api";
import type { Organization } from "@/features/dashboard/types/organization.types";

const SELECTED_ORG_KEY = "selected_org";

interface OrgContextValue {
  organizations: Organization[];
  currentOrg: Organization | null;
  switchOrg: (org: Organization) => void;
  addOrg: (org: Organization) => void;
  removeOrg: (orgId: string) => void;
}

const OrgContext = createContext<OrgContextValue | null>(null);

export function OrgProvider({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);

  useEffect(() => {
    if (!isAuthenticated) {
      setOrganizations([]);
      setCurrentOrg(null);
      return;
    }

    async function load() {
      try {
        const orgs = await fetchOrganizations();
        setOrganizations(orgs);

        const storedId = localStorage.getItem(SELECTED_ORG_KEY);
        const matched = orgs.find((o) => o.id === storedId);
        if (matched) {
          setCurrentOrg(matched);
        } else if (orgs.length > 0) {
          setCurrentOrg(orgs[0]);
          localStorage.setItem(SELECTED_ORG_KEY, orgs[0].id);
        }
      } catch {
        // silently fail — user may not have orgs yet
      }
    }

    load();
  }, [isAuthenticated]);

  const switchOrg = useCallback((org: Organization) => {
    setCurrentOrg(org);
    localStorage.setItem(SELECTED_ORG_KEY, org.id);
  }, []);

  const addOrg = useCallback((org: Organization) => {
    setOrganizations((prev) => [...prev, org]);
    setCurrentOrg(org);
    localStorage.setItem(SELECTED_ORG_KEY, org.id);
  }, []);

  const removeOrg = useCallback((orgId: string) => {
    setOrganizations((prev) => {
      const remaining = prev.filter((o) => o.id !== orgId);
      if (remaining.length > 0) {
        setCurrentOrg(remaining[0]);
        localStorage.setItem(SELECTED_ORG_KEY, remaining[0].id);
      } else {
        setCurrentOrg(null);
        localStorage.removeItem(SELECTED_ORG_KEY);
      }
      return remaining;
    });
  }, []);

  const value = useMemo<OrgContextValue>(
    () => ({ organizations, currentOrg, switchOrg, addOrg, removeOrg }),
    [organizations, currentOrg, switchOrg, addOrg, removeOrg]
  );

  return <OrgContext.Provider value={value}>{children}</OrgContext.Provider>;
}

export function useOrg(): OrgContextValue {
  const context = useContext(OrgContext);
  if (!context) {
    throw new Error("useOrg must be used within an OrgProvider");
  }
  return context;
}
