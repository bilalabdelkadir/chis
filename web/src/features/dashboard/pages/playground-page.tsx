import { useState } from "react";
import { Link } from "react-router";
import { Send, Loader2, Copy, Check, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PageHeader } from "../components/page-header";
import { sendTestWebhook } from "../api/dashboard-api";
import type { SendWebhookResponse } from "../types/dashboard.types";

const HTTP_METHODS = ["POST", "GET", "PUT", "DELETE", "PATCH"] as const;

const DEFAULT_PAYLOAD = JSON.stringify(
  {
    event: "test",
    data: {
      message: "Hello from Chis!",
    },
  },
  null,
  2,
);

export function PlaygroundPage() {
  const [apiKey, setApiKey] = useState("");
  const [method, setMethod] = useState("POST");
  const [url, setUrl] = useState("");
  const [payload, setPayload] = useState(DEFAULT_PAYLOAD);

  const [isLoading, setIsLoading] = useState(false);
  const [response, setResponse] = useState<SendWebhookResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  function handleFormatJson() {
    try {
      const parsed = JSON.parse(payload);
      setPayload(JSON.stringify(parsed, null, 2));
    } catch {
      // invalid JSON — leave as-is
    }
  }

  async function handleSend() {
    setError(null);
    setResponse(null);

    if (!apiKey.trim()) {
      setError("API key is required.");
      return;
    }
    if (!url.trim()) {
      setError("URL is required.");
      return;
    }

    let parsedPayload: unknown;
    try {
      parsedPayload = JSON.parse(payload);
    } catch {
      setError("Invalid JSON payload.");
      return;
    }

    setIsLoading(true);
    try {
      const result = await sendTestWebhook(apiKey.trim(), {
        url: url.trim(),
        method,
        payload: parsedPayload,
      });
      setResponse(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An unexpected error occurred.");
    } finally {
      setIsLoading(false);
    }
  }

  function handleCopyId() {
    if (!response) return;
    navigator.clipboard.writeText(response.messageId);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <>
      <PageHeader
        title="Playground"
        description="Send test webhooks to any URL. Useful for debugging and integration checks."
      />

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Left panel — Request Builder */}
        <Card className="border-border/80">
          <CardHeader>
            <CardTitle className="text-sm font-medium">Request</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="api-key">API Key</Label>
              <Input
                id="api-key"
                type="password"
                placeholder="chis_sk_..."
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
              />
              <p className="text-muted-foreground text-xs">
                Paste your full API key.{" "}
                <Link
                  to="/dashboard/api-keys"
                  className="text-primary underline underline-offset-4"
                >
                  Create one
                </Link>{" "}
                if you don't have one yet.
              </p>
            </div>

            <div className="grid grid-cols-[120px_1fr] gap-3">
              <div className="space-y-1.5">
                <Label>Method</Label>
                <Select value={method} onValueChange={(v) => v && setMethod(v)}>
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {HTTP_METHODS.map((m) => (
                      <SelectItem key={m} value={m}>
                        {m}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="url">URL</Label>
                <Input
                  id="url"
                  placeholder="https://httpbin.org/post"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <Label htmlFor="payload">Payload</Label>
                <Button
                  variant="ghost"
                  size="xs"
                  onClick={handleFormatJson}
                >
                  Format JSON
                </Button>
              </div>
              <Textarea
                id="payload"
                value={payload}
                onChange={(e) => setPayload(e.target.value)}
                className="min-h-40 font-mono text-xs"
              />
            </div>

            <Button
              className="w-full"
              onClick={handleSend}
              disabled={isLoading}
            >
              {isLoading ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <Send className="size-4" />
              )}
              {isLoading ? "Sending..." : "Send Webhook"}
            </Button>
          </CardContent>
        </Card>

        {/* Right panel — Response */}
        <div className="space-y-6">
          <Card className="border-border/80">
            <CardHeader>
              <CardTitle className="text-sm font-medium">Response</CardTitle>
            </CardHeader>
            <CardContent>
              {error && (
                <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
                  {error}
                </div>
              )}

              {response && (
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">{response.status}</Badge>
                  </div>
                  <div className="space-y-1">
                    <p className="text-muted-foreground text-xs font-medium">
                      Message ID
                    </p>
                    <div className="flex items-center gap-2">
                      <code className="bg-muted text-xs px-2 py-1 rounded break-all flex-1">
                        {response.messageId}
                      </code>
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        onClick={handleCopyId}
                      >
                        {copied ? (
                          <Check className="size-3.5" />
                        ) : (
                          <Copy className="size-3.5" />
                        )}
                      </Button>
                    </div>
                  </div>
                  <Link
                    to="/dashboard/logs"
                    className="inline-flex items-center gap-1.5 text-xs text-primary underline underline-offset-4"
                  >
                    View in logs
                    <ExternalLink className="size-3" />
                  </Link>
                </div>
              )}

              {!error && !response && (
                <p className="text-muted-foreground text-sm">
                  Send a request to see the response here.
                </p>
              )}
            </CardContent>
          </Card>

          <Card className="border-border/80">
            <CardHeader>
              <CardTitle className="text-sm font-medium">
                How it works
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm text-muted-foreground">
              <p>
                Sends a webhook to any URL. The request is queued and delivered
                asynchronously with automatic retries on failure.
              </p>
              <div className="space-y-1">
                <p className="text-foreground text-xs font-medium">
                  API Endpoint
                </p>
                <code className="bg-muted text-xs px-2 py-1 rounded block">
                  POST /webhook/send
                </code>
              </div>
              <div className="space-y-1">
                <p className="text-foreground text-xs font-medium">Headers</p>
                <code className="bg-muted text-xs px-2 py-1 rounded block">
                  X-API-Key: your_api_key
                </code>
              </div>
              <div className="space-y-1">
                <p className="text-foreground text-xs font-medium">
                  Example Payload
                </p>
                <pre className="bg-muted text-xs px-2 py-1.5 rounded overflow-x-auto">
{`{
  "url": "https://example.com/hook",
  "method": "POST",
  "payload": { "event": "test" }
}`}
                </pre>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
