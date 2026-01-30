import { useState } from "react";
import { Building2, Check, ChevronDown, Plus } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOrg } from "@/shared/context/org-context";
import { createOrganization } from "../api/organization-api";

export function OrgSwitcher() {
  const { organizations, currentOrg, switchOrg, addOrg } = useOrg();
  const [newOrgName, setNewOrgName] = useState("");
  const [isCreating, setIsCreating] = useState(false);

  async function handleCreate() {
    const name = newOrgName.trim();
    if (!name) return;

    setIsCreating(true);
    try {
      const org = await createOrganization({ name });
      addOrg(org);
      setNewOrgName("");
    } catch {
      // silently fail
    } finally {
      setIsCreating(false);
    }
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className="flex items-center gap-2 rounded-md px-3 py-1.5 text-sm font-medium hover:bg-accent focus:outline-none"
      >
        <Building2 className="size-4" />
        <span>{currentOrg?.name ?? "Select org"}</span>
        <ChevronDown className="size-3 text-muted-foreground" />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" sideOffset={8}>
        <DropdownMenuLabel>Organizations</DropdownMenuLabel>
        {organizations.map((org) => (
          <DropdownMenuItem
            key={org.id}
            onClick={() => switchOrg(org)}
            className="flex items-center justify-between"
          >
            <span>{org.name}</span>
            {org.id === currentOrg?.id && (
              <Check className="size-4 text-primary" />
            )}
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator />
        <div className="px-1.5 py-1.5">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleCreate();
            }}
            className="flex items-center gap-1.5"
          >
            <Input
              value={newOrgName}
              onChange={(e) => setNewOrgName(e.target.value)}
              placeholder="New organization"
              className="h-7 text-xs"
              onClick={(e) => e.stopPropagation()}
              onKeyDown={(e) => e.stopPropagation()}
            />
            <Button
              type="submit"
              size="sm"
              variant="ghost"
              disabled={isCreating || !newOrgName.trim()}
              className="h-7 px-2"
            >
              <Plus className="size-3" />
            </Button>
          </form>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
