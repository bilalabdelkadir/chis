import { useState } from "react";

const tabs = [
  {
    id: "typescript",
    label: "TypeScript",
    code: `import crypto from "crypto";
import type { Request, Response, NextFunction } from "express";

export function verifyWebhook(secret: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    const msgId = req.headers["x-webhook-id"] as string;
    const timestamp = req.headers["x-webhook-timestamp"] as string;
    const signature = req.headers["x-webhook-signature"] as string;
    if (!msgId || !timestamp || !signature) {
      return res.status(401).json({ error: "Missing webhook headers" });
    }

    const now = Math.floor(Date.now() / 1000);
    if (Math.abs(now - parseInt(timestamp)) > 300) {
      return res.status(401).json({ error: "Timestamp too old" });
    }

    const body = JSON.stringify(req.body);
    const signedContent = \`\${msgId}.\${timestamp}.\${body}\`;
    const secretBytes = Buffer.from(secret.replace("whsec_", ""), "base64");
    const expected = crypto
      .createHmac("sha256", secretBytes)
      .update(signedContent)
      .digest("base64");

    const expectedSig = \`v1,\${expected}\`;
    const signatures = signature.split(" ");
    const valid = signatures.some(
      (sig) =>
        sig.length === expectedSig.length &&
        crypto.timingSafeEqual(Buffer.from(sig), Buffer.from(expectedSig))
    );

    if (!valid) return res.status(401).json({ error: "Invalid signature" });
    next();
  };
}`,
  },
  {
    id: "nodejs",
    label: "Node.js",
    code: `const crypto = require("crypto");
const http = require("http");

function verifyWebhook(req, body, secret) {
  const msgId = req.headers["x-webhook-id"];
  const timestamp = req.headers["x-webhook-timestamp"];
  const signature = req.headers["x-webhook-signature"];
  if (!msgId || !timestamp || !signature) return false;

  const now = Math.floor(Date.now() / 1000);
  if (Math.abs(now - parseInt(timestamp)) > 300) return false;

  const signedContent = \`\${msgId}.\${timestamp}.\${body}\`;
  const secretBytes = Buffer.from(secret.replace("whsec_", ""), "base64");
  const expected = crypto
    .createHmac("sha256", secretBytes)
    .update(signedContent)
    .digest("base64");

  const expectedSig = \`v1,\${expected}\`;
  return signature.split(" ").some(
    (sig) =>
      sig.length === expectedSig.length &&
      crypto.timingSafeEqual(Buffer.from(sig), Buffer.from(expectedSig))
  );
}`,
  },
  {
    id: "python-fastapi",
    label: "Python (FastAPI)",
    code: `import hmac, hashlib, base64, time
from fastapi import Request, HTTPException

async def verify_webhook(request: Request, secret: str):
    msg_id = request.headers.get("x-webhook-id")
    timestamp = request.headers.get("x-webhook-timestamp")
    signature = request.headers.get("x-webhook-signature")
    if not all([msg_id, timestamp, signature]):
        raise HTTPException(401, "Missing webhook headers")

    if abs(time.time() - int(timestamp)) > 300:
        raise HTTPException(401, "Timestamp too old")

    body = (await request.body()).decode()
    signed_content = f"{msg_id}.{timestamp}.{body}"
    secret_bytes = base64.b64decode(secret.removeprefix("whsec_"))
    expected = base64.b64encode(
        hmac.new(secret_bytes, signed_content.encode(), hashlib.sha256).digest()
    ).decode()

    expected_sig = f"v1,{expected}"
    if not any(
        hmac.compare_digest(sig, expected_sig)
        for sig in signature.split(" ")
    ):
        raise HTTPException(401, "Invalid signature")`,
  },
  {
    id: "python-django",
    label: "Python (Django)",
    code: `import hmac, hashlib, base64, time
from django.http import JsonResponse

def verify_webhook(view_func):
    def wrapper(request, *args, **kwargs):
        secret = settings.WEBHOOK_SECRET
        msg_id = request.headers.get("X-Webhook-Id")
        timestamp = request.headers.get("X-Webhook-Timestamp")
        signature = request.headers.get("X-Webhook-Signature")
        if not all([msg_id, timestamp, signature]):
            return JsonResponse({"error": "Missing headers"}, status=401)

        if abs(time.time() - int(timestamp)) > 300:
            return JsonResponse({"error": "Timestamp too old"}, status=401)

        body = request.body.decode()
        signed_content = f"{msg_id}.{timestamp}.{body}"
        secret_bytes = base64.b64decode(secret.removeprefix("whsec_"))
        expected = base64.b64encode(
            hmac.new(secret_bytes, signed_content.encode(), hashlib.sha256).digest()
        ).decode()

        expected_sig = f"v1,{expected}"
        if not any(
            hmac.compare_digest(sig, expected_sig)
            for sig in signature.split(" ")
        ):
            return JsonResponse({"error": "Invalid signature"}, status=401)
        return view_func(request, *args, **kwargs)
    return wrapper`,
  },
  {
    id: "go",
    label: "Go",
    code: `package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

func verifyWebhook(r *http.Request, body []byte, secret string) error {
	msgId := r.Header.Get("X-Webhook-Id")
	timestamp := r.Header.Get("X-Webhook-Timestamp")
	signature := r.Header.Get("X-Webhook-Signature")
	if msgId == "" || timestamp == "" || signature == "" {
		return fmt.Errorf("missing webhook headers")
	}

	ts, _ := strconv.ParseInt(timestamp, 10, 64)
	if math.Abs(float64(time.Now().Unix()-ts)) > 300 {
		return fmt.Errorf("timestamp too old")
	}

	signedContent := fmt.Sprintf("%s.%s.%s", msgId, timestamp, string(body))
	secretBytes, _ := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(secret, "whsec_"),
	)
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signedContent))
	expected := "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	for _, sig := range strings.Split(signature, " ") {
		if hmac.Equal([]byte(sig), []byte(expected)) {
			return nil
		}
	}
	return fmt.Errorf("invalid signature")
}`,
  },
  {
    id: "java",
    label: "Java (Spring)",
    code: `import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.util.Base64;

@RestController
public class WebhookController {

    @PostMapping("/webhook")
    public ResponseEntity<String> handleWebhook(
            @RequestHeader("X-Webhook-Id") String msgId,
            @RequestHeader("X-Webhook-Timestamp") String timestamp,
            @RequestHeader("X-Webhook-Signature") String signature,
            @RequestBody String body) throws Exception {

        long now = System.currentTimeMillis() / 1000;
        if (Math.abs(now - Long.parseLong(timestamp)) > 300) {
            return ResponseEntity.status(401).body("Timestamp too old");
        }

        String secret = System.getenv("WEBHOOK_SECRET");
        byte[] secretBytes = Base64.getDecoder().decode(
            secret.replaceFirst("^whsec_", "")
        );
        String signedContent = msgId + "." + timestamp + "." + body;

        Mac mac = Mac.getInstance("HmacSHA256");
        mac.init(new SecretKeySpec(secretBytes, "HmacSHA256"));
        String expected = "v1," + Base64.getEncoder().encodeToString(
            mac.doFinal(signedContent.getBytes())
        );

        boolean valid = false;
        for (String sig : signature.split(" ")) {
            if (MessageDigest.isEqual(sig.getBytes(), expected.getBytes())) {
                valid = true;
                break;
            }
        }
        if (!valid) return ResponseEntity.status(401).body("Invalid signature");
        return ResponseEntity.ok("OK");
    }
}`,
  },
  {
    id: "php",
    label: "PHP (Laravel)",
    code: `<?php

namespace App\\Http\\Middleware;

use Closure;
use Illuminate\\Http\\Request;

class VerifyWebhook
{
    public function handle(Request $request, Closure $next)
    {
        $msgId = $request->header('X-Webhook-Id');
        $timestamp = $request->header('X-Webhook-Timestamp');
        $signature = $request->header('X-Webhook-Signature');
        if (!$msgId || !$timestamp || !$signature) {
            return response()->json(['error' => 'Missing headers'], 401);
        }

        if (abs(time() - intval($timestamp)) > 300) {
            return response()->json(['error' => 'Timestamp too old'], 401);
        }

        $body = $request->getContent();
        $signedContent = "{$msgId}.{$timestamp}.{$body}";
        $secret = config('services.webhook.secret');
        $secretBytes = base64_decode(
            str_replace('whsec_', '', $secret)
        );
        $expected = 'v1,' . base64_encode(
            hash_hmac('sha256', $signedContent, $secretBytes, true)
        );

        $valid = collect(explode(' ', $signature))
            ->contains(fn ($sig) => hash_equals($expected, $sig));

        if (!$valid) {
            return response()->json(['error' => 'Invalid signature'], 401);
        }
        return $next($request);
    }
}`,
  },
];

export default function CodeTabs() {
  const [activeTab, setActiveTab] = useState("typescript");
  const [copied, setCopied] = useState(false);

  const activeSnippet = tabs.find((t) => t.id === activeTab)!;

  const handleCopy = async () => {
    await navigator.clipboard.writeText(activeSnippet.code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="rounded-xl border border-surface-800 bg-surface-900 overflow-hidden">
      <div className="flex items-center border-b border-surface-800">
        <div className="flex-1 flex overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 text-sm font-medium whitespace-nowrap transition-colors ${
                activeTab === tab.id
                  ? "text-brand-400 border-b-2 border-brand-400"
                  : "text-surface-400 hover:text-surface-300"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <button
          onClick={handleCopy}
          className="flex-shrink-0 px-4 py-3 text-xs font-medium text-surface-400 hover:text-surface-200 transition-colors"
        >
          {copied ? "Copied!" : "Copy"}
        </button>
      </div>
      <pre className="p-5 text-sm leading-relaxed overflow-x-auto font-mono">
        <code className="text-surface-300">{activeSnippet.code}</code>
      </pre>
    </div>
  );
}
