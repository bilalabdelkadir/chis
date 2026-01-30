import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { acceptInvitation } from "@/features/dashboard/api/invitation-api";
import { ApiRequestError } from "@/shared/api/api-error";
import { useOrg } from "@/shared/context/org-context";
import type { AcceptInvitationResponse } from "@/features/dashboard/types/invitation.types";

export function AcceptInvitePage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");

  const { addOrg } = useOrg();
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<AcceptInvitationResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!token) return;

    let cancelled = false;

    async function accept() {
      setIsLoading(true);
      setError(null);

      try {
        const data = await acceptInvitation(token!);
        if (!cancelled) {
          setResult(data);
          addOrg({
            id: data.orgId,
            name: data.orgName,
            slug: data.orgSlug,
            role: data.role,
          });
        }
      } catch (err) {
        if (!cancelled) {
          if (err instanceof ApiRequestError) {
            setError(err.message);
          } else {
            setError("An unexpected error occurred");
          }
        }
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    accept();
    return () => {
      cancelled = true;
    };
  }, [token]);

  if (!token) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-md border-border">
          <CardHeader className="px-6 pt-6">
            <CardTitle className="text-xl">Invalid invitation</CardTitle>
            <CardDescription>
              This invitation link is invalid or missing a token.
            </CardDescription>
          </CardHeader>
          <CardContent className="px-6 pb-6">
            <Button render={<Link to="/dashboard" />} className="w-full">
              Go to dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-md border-border">
          <CardHeader className="px-6 pt-6">
            <CardTitle className="text-xl">Accepting invitation...</CardTitle>
            <CardDescription>Please wait while we process your invitation.</CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-md border-border">
          <CardHeader className="px-6 pt-6">
            <CardTitle className="text-xl">Invitation failed</CardTitle>
            <CardDescription>{error}</CardDescription>
          </CardHeader>
          <CardContent className="px-6 pb-6">
            <Button render={<Link to="/dashboard" />} className="w-full">
              Go to dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (result) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-md border-border">
          <CardHeader className="px-6 pt-6">
            <CardTitle className="text-xl">You're in!</CardTitle>
            <CardDescription>
              You've joined <strong>{result.orgName}</strong> as a {result.role}.
            </CardDescription>
          </CardHeader>
          <CardContent className="px-6 pb-6">
            <Button render={<Link to="/dashboard" />} className="w-full">
              Go to dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return null;
}
