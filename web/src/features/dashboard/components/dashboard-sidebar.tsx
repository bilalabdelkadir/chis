import { NavLink, useNavigate } from 'react-router';
import {
  BarChart3,
  Key,
  ScrollText,
  Play,
  Users,
  Settings,
  LogOut,
  ChevronDown,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAuth } from '@/shared/hooks/use-auth';
import { useOrg } from '@/shared/context/org-context';
import { OrgSwitcher } from './org-switcher';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Separator } from '@/components/ui/separator';

const mainNavItems = [
  { to: '/dashboard', label: 'Overview', icon: BarChart3, end: true },
  { to: '/dashboard/api-keys', label: 'API Keys', icon: Key, end: false },
  { to: '/dashboard/logs', label: 'Logs', icon: ScrollText, end: false },
  { to: '/dashboard/playground', label: 'Playground', icon: Play, end: false },
] as const;

const adminNavItems = [
  { to: '/dashboard/members', label: 'Members', icon: Users, end: false },
  { to: '/dashboard/settings', label: 'Settings', icon: Settings, end: false },
] as const;

export function DashboardSidebar({ onNavigate }: { onNavigate?: () => void }) {
  const { user, logout } = useAuth();
  const { currentOrg } = useOrg();
  const isAdmin = currentOrg?.role === 'admin';
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate('/login', { replace: true });
  }

  const navLinkClass = ({ isActive }: { isActive: boolean }) =>
    cn(
      'flex items-center gap-3 rounded-lg px-2.5 py-2 text-sm font-medium transition-colors',
      isActive
        ? 'bg-sidebar-accent text-sidebar-accent-foreground'
        : 'text-sidebar-foreground/80 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground',
    );

  return (
    <aside className="flex h-full w-[240px] shrink-0 flex-col border-r border-sidebar-border bg-sidebar">
      {/* Brand */}
      <div className="flex h-14 shrink-0 items-center gap-2 border-b border-sidebar-border px-4">
        <NavLink
          to="/dashboard"
          onClick={onNavigate}
          className="flex items-center gap-2 font-semibold text-sidebar-foreground no-underline"
        >
          <img src="/logo.png" alt="Chis" className="h-8 w-auto dark:invert" />
        </NavLink>
      </div>

      {/* Org switcher */}
      <div className="shrink-0 border-b border-sidebar-border px-3 py-3">
        <OrgSwitcher variant="sidebar" />
      </div>

      {/* Main nav */}
      <nav className="flex-1 overflow-y-auto px-3 py-3">
        <ul className="space-y-0.5">
          {mainNavItems.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                end={item.end}
                onClick={onNavigate}
                className={navLinkClass}
              >
                <item.icon className="size-4 shrink-0 opacity-80" />
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>

        {isAdmin && (
          <>
            <Separator className="my-3 bg-sidebar-border" />
            <p className="mb-1.5 px-2.5 text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider">
              Organization
            </p>
            <ul className="space-y-0.5">
              {adminNavItems.map((item) => (
                <li key={item.to}>
                  <NavLink
                    to={item.to}
                    end={item.end}
                    onClick={onNavigate}
                    className={navLinkClass}
                  >
                    <item.icon className="size-4 shrink-0 opacity-80" />
                    {item.label}
                  </NavLink>
                </li>
              ))}
            </ul>
          </>
        )}
      </nav>

      {/* User menu */}
      <div className="shrink-0 border-t border-sidebar-border p-3">
        <DropdownMenu>
          <DropdownMenuTrigger>
            <Button
              variant="ghost"
              className="w-full justify-between gap-2 rounded-lg px-2.5 py-2 h-auto font-normal text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
            >
              <span className="min-w-0 truncate text-left">
                <span className="block truncate text-sm font-medium">
                  {user?.firstName || 'Account'}
                </span>
                <span className="block truncate text-xs text-sidebar-foreground/70">
                  {user?.email}
                </span>
              </span>
              <ChevronDown className="size-4 shrink-0 opacity-70" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" side="top" className="w-56">
            <DropdownMenuItem onClick={handleLogout} className="gap-2">
              <LogOut className="size-4" />
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>
  );
}
