import { useState } from "react";

const c = {
  kw: "text-brand-400",
  str: "text-green-400",
  fn: "text-yellow-400",
  cm: "text-surface-600",
  num: "text-orange-400",
  dim: "text-surface-400",
  key: "text-sky-400",
  plain: "text-surface-300",
};

type Line = { text: string; cls?: string }[];

const tabs: { id: string; label: string; lines: Line[] }[] = [
  {
    id: "nodejs",
    label: "Node.js",
    lines: [
      [
        { text: "const ", cls: c.kw },
        { text: "crypto ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "require", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"crypto"', cls: c.str },
        { text: ");", cls: c.dim },
      ],
      [],
      [
        { text: "function ", cls: c.kw },
        { text: "verify", cls: c.fn },
        { text: "(body, secret, headers) {", cls: c.dim },
      ],
      [
        { text: "  const ", cls: c.kw },
        { text: "msgId ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-id"', cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [
        { text: "  const ", cls: c.kw },
        { text: "ts    ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-timestamp"', cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [
        { text: "  const ", cls: c.kw },
        { text: "sig   ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-signature"', cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [],
      [
        { text: "  // Reject if older than 5 minutes", cls: c.cm },
      ],
      [
        { text: "  if ", cls: c.kw },
        { text: "(Math.", cls: c.plain },
        { text: "abs", cls: c.fn },
        { text: "(Date.", cls: c.plain },
        { text: "now", cls: c.fn },
        { text: "() / ", cls: c.dim },
        { text: "1000", cls: c.num },
        { text: " - ts) > ", cls: c.dim },
        { text: "300", cls: c.num },
        { text: ") return ", cls: c.dim },
        { text: "false", cls: c.kw },
        { text: ";", cls: c.dim },
      ],
      [],
      [
        { text: "  const ", cls: c.kw },
        { text: "signed ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "`${msgId}.${ts}.${body}`", cls: c.str },
        { text: ";", cls: c.dim },
      ],
      [
        { text: "  const ", cls: c.kw },
        { text: "key ", cls: c.plain },
        { text: "= Buffer.", cls: c.plain },
        { text: "from", cls: c.fn },
        { text: "(secret.", cls: c.dim },
        { text: "replace", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"whsec_"', cls: c.str },
        { text: ", ", cls: c.dim },
        { text: '""', cls: c.str },
        { text: "), ", cls: c.dim },
        { text: '"base64"', cls: c.str },
        { text: ");", cls: c.dim },
      ],
      [
        { text: "  const ", cls: c.kw },
        { text: "expected ", cls: c.plain },
        { text: "= crypto.", cls: c.plain },
        { text: "createHmac", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"sha256"', cls: c.str },
        { text: ", key).", cls: c.dim },
        { text: "update", cls: c.fn },
        { text: "(signed).", cls: c.dim },
        { text: "digest", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"base64"', cls: c.str },
        { text: ");", cls: c.dim },
      ],
      [],
      [
        { text: "  return ", cls: c.kw },
        { text: "sig.", cls: c.plain },
        { text: "split", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '" "', cls: c.str },
        { text: ").", cls: c.dim },
        { text: "some", cls: c.fn },
        { text: "(s =>", cls: c.dim },
      ],
      [
        { text: "    s === ", cls: c.dim },
        { text: "`v1,${expected}`", cls: c.str },
      ],
      [
        { text: "  );", cls: c.dim },
      ],
      [
        { text: "}", cls: c.dim },
      ],
    ],
  },
  {
    id: "python",
    label: "Python",
    lines: [
      [
        { text: "import ", cls: c.kw },
        { text: "hmac, hashlib, base64, time", cls: c.plain },
      ],
      [],
      [
        { text: "def ", cls: c.kw },
        { text: "verify", cls: c.fn },
        { text: "(body, secret, headers):", cls: c.dim },
      ],
      [
        { text: "    msg_id ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-id"', cls: c.str },
        { text: "]", cls: c.dim },
      ],
      [
        { text: "    ts     ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-timestamp"', cls: c.str },
        { text: "]", cls: c.dim },
      ],
      [
        { text: "    sig    ", cls: c.plain },
        { text: "= headers[", cls: c.dim },
        { text: '"x-webhook-signature"', cls: c.str },
        { text: "]", cls: c.dim },
      ],
      [],
      [
        { text: "    if ", cls: c.kw },
        { text: "abs", cls: c.fn },
        { text: "(time.", cls: c.plain },
        { text: "time", cls: c.fn },
        { text: "() - ", cls: c.dim },
        { text: "int", cls: c.fn },
        { text: "(ts)) > ", cls: c.dim },
        { text: "300", cls: c.num },
        { text: ":", cls: c.dim },
      ],
      [
        { text: "        return ", cls: c.kw },
        { text: "False", cls: c.kw },
      ],
      [],
      [
        { text: "    key ", cls: c.plain },
        { text: "= base64.", cls: c.plain },
        { text: "b64decode", cls: c.fn },
        { text: "(secret.", cls: c.dim },
        { text: "removeprefix", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"whsec_"', cls: c.str },
        { text: "))", cls: c.dim },
      ],
      [
        { text: "    signed ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: 'f"{msg_id}.{ts}.{body}"', cls: c.str },
      ],
      [
        { text: "    expected ", cls: c.plain },
        { text: "= base64.", cls: c.plain },
        { text: "b64encode", cls: c.fn },
        { text: "(", cls: c.dim },
      ],
      [
        { text: "        hmac.", cls: c.plain },
        { text: "new", cls: c.fn },
        { text: "(key, signed.", cls: c.dim },
        { text: "encode", cls: c.fn },
        { text: "(), hashlib.", cls: c.dim },
        { text: "sha256", cls: c.key },
        { text: ").", cls: c.dim },
        { text: "digest", cls: c.fn },
        { text: "()", cls: c.dim },
      ],
      [
        { text: "    ).", cls: c.dim },
        { text: "decode", cls: c.fn },
        { text: "()", cls: c.dim },
      ],
      [],
      [
        { text: "    return ", cls: c.kw },
        { text: "any", cls: c.fn },
        { text: "(hmac.", cls: c.dim },
        { text: "compare_digest", cls: c.fn },
        { text: "(s, ", cls: c.dim },
        { text: 'f"v1,{expected}"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "               for ", cls: c.kw },
        { text: "s ", cls: c.plain },
        { text: "in ", cls: c.kw },
        { text: "sig.", cls: c.dim },
        { text: "split", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '" "', cls: c.str },
        { text: "))", cls: c.dim },
      ],
    ],
  },
  {
    id: "go",
    label: "Go",
    lines: [
      [
        { text: "func ", cls: c.kw },
        { text: "verify", cls: c.fn },
        { text: "(body []", cls: c.dim },
        { text: "byte", cls: c.kw },
        { text: ", secret ", cls: c.dim },
        { text: "string", cls: c.kw },
        { text: ", h http.", cls: c.dim },
        { text: "Header", cls: c.fn },
        { text: ") ", cls: c.dim },
        { text: "bool", cls: c.kw },
        { text: " {", cls: c.dim },
      ],
      [
        { text: "    msgId ", cls: c.plain },
        { text: ":= h.", cls: c.dim },
        { text: "Get", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"X-Webhook-Id"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "    ts    ", cls: c.plain },
        { text: ":= h.", cls: c.dim },
        { text: "Get", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"X-Webhook-Timestamp"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "    sig   ", cls: c.plain },
        { text: ":= h.", cls: c.dim },
        { text: "Get", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"X-Webhook-Signature"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [],
      [
        { text: "    t, _ ", cls: c.plain },
        { text: ":= strconv.", cls: c.dim },
        { text: "ParseInt", cls: c.fn },
        { text: "(ts, ", cls: c.dim },
        { text: "10", cls: c.num },
        { text: ", ", cls: c.dim },
        { text: "64", cls: c.num },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "    if ", cls: c.kw },
        { text: "math.", cls: c.plain },
        { text: "Abs", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "float64", cls: c.fn },
        { text: "(time.", cls: c.dim },
        { text: "Now", cls: c.fn },
        { text: "().", cls: c.dim },
        { text: "Unix", cls: c.fn },
        { text: "() - t)) > ", cls: c.dim },
        { text: "300", cls: c.num },
        { text: " {", cls: c.dim },
      ],
      [
        { text: "        return ", cls: c.kw },
        { text: "false", cls: c.kw },
      ],
      [
        { text: "    }", cls: c.dim },
      ],
      [],
      [
        { text: "    key, _ ", cls: c.plain },
        { text: ":= base64.", cls: c.dim },
        { text: "StdEncoding.", cls: c.plain },
        { text: "DecodeString", cls: c.fn },
        { text: "(", cls: c.dim },
      ],
      [
        { text: "        strings.", cls: c.plain },
        { text: "TrimPrefix", cls: c.fn },
        { text: "(secret, ", cls: c.dim },
        { text: '"whsec_"', cls: c.str },
        { text: "))", cls: c.dim },
      ],
      [
        { text: "    mac ", cls: c.plain },
        { text: ":= hmac.", cls: c.dim },
        { text: "New", cls: c.fn },
        { text: "(sha256.", cls: c.dim },
        { text: "New", cls: c.fn },
        { text: ", key)", cls: c.dim },
      ],
      [
        { text: "    mac.", cls: c.plain },
        { text: "Write", cls: c.fn },
        { text: "([]", cls: c.dim },
        { text: "byte", cls: c.kw },
        { text: "(msgId + ", cls: c.dim },
        { text: '"."', cls: c.str },
        { text: " + ts + ", cls: c.dim },
        { text: '"."', cls: c.str },
        { text: " + ", cls: c.dim },
        { text: "string", cls: c.kw },
        { text: "(body)))", cls: c.dim },
      ],
      [
        { text: "    expected ", cls: c.plain },
        { text: ":= ", cls: c.dim },
        { text: '"v1,"', cls: c.str },
        { text: " + base64.", cls: c.dim },
        { text: "StdEncoding.", cls: c.plain },
        { text: "EncodeToString", cls: c.fn },
        { text: "(mac.", cls: c.dim },
        { text: "Sum", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "nil", cls: c.kw },
        { text: "))", cls: c.dim },
      ],
      [],
      [
        { text: "    for ", cls: c.kw },
        { text: "_, s ", cls: c.plain },
        { text: ":= ", cls: c.dim },
        { text: "range ", cls: c.kw },
        { text: "strings.", cls: c.plain },
        { text: "Split", cls: c.fn },
        { text: "(sig, ", cls: c.dim },
        { text: '" "', cls: c.str },
        { text: ") {", cls: c.dim },
      ],
      [
        { text: "        if ", cls: c.kw },
        { text: "hmac.", cls: c.plain },
        { text: "Equal", cls: c.fn },
        { text: "([]", cls: c.dim },
        { text: "byte", cls: c.kw },
        { text: "(s), []", cls: c.dim },
        { text: "byte", cls: c.kw },
        { text: "(expected)) { ", cls: c.dim },
        { text: "return ", cls: c.kw },
        { text: "true", cls: c.kw },
        { text: " }", cls: c.dim },
      ],
      [
        { text: "    }", cls: c.dim },
      ],
      [
        { text: "    return ", cls: c.kw },
        { text: "false", cls: c.kw },
      ],
      [
        { text: "}", cls: c.dim },
      ],
    ],
  },
  {
    id: "php",
    label: "PHP",
    lines: [
      [
        { text: "function ", cls: c.kw },
        { text: "verify", cls: c.fn },
        { text: "($body, $secret, $headers) {", cls: c.dim },
      ],
      [
        { text: "    $msgId ", cls: c.plain },
        { text: "= $headers[", cls: c.dim },
        { text: "'X-Webhook-Id'", cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [
        { text: "    $ts    ", cls: c.plain },
        { text: "= $headers[", cls: c.dim },
        { text: "'X-Webhook-Timestamp'", cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [
        { text: "    $sig   ", cls: c.plain },
        { text: "= $headers[", cls: c.dim },
        { text: "'X-Webhook-Signature'", cls: c.str },
        { text: "];", cls: c.dim },
      ],
      [],
      [
        { text: "    if ", cls: c.kw },
        { text: "(", cls: c.dim },
        { text: "abs", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "time", cls: c.fn },
        { text: "() - (", cls: c.dim },
        { text: "int", cls: c.kw },
        { text: ")$ts) > ", cls: c.dim },
        { text: "300", cls: c.num },
        { text: ") return ", cls: c.dim },
        { text: "false", cls: c.kw },
        { text: ";", cls: c.dim },
      ],
      [],
      [
        { text: "    $key ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "base64_decode", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "str_replace", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "'whsec_'", cls: c.str },
        { text: ", ", cls: c.dim },
        { text: "''", cls: c.str },
        { text: ", $secret));", cls: c.dim },
      ],
      [
        { text: "    $signed ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: '"$msgId.$ts.$body"', cls: c.str },
        { text: ";", cls: c.dim },
      ],
      [
        { text: "    $expected ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "'v1,' ", cls: c.str },
        { text: ". ", cls: c.dim },
        { text: "base64_encode", cls: c.fn },
        { text: "(", cls: c.dim },
      ],
      [
        { text: "        ", cls: c.dim },
        { text: "hash_hmac", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "'sha256'", cls: c.str },
        { text: ", $signed, $key, ", cls: c.dim },
        { text: "true", cls: c.kw },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "    );", cls: c.dim },
      ],
      [],
      [
        { text: "    foreach ", cls: c.kw },
        { text: "(", cls: c.dim },
        { text: "explode", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "' '", cls: c.str },
        { text: ", $sig) ", cls: c.dim },
        { text: "as ", cls: c.kw },
        { text: "$s) {", cls: c.dim },
      ],
      [
        { text: "        if ", cls: c.kw },
        { text: "(", cls: c.dim },
        { text: "hash_equals", cls: c.fn },
        { text: "($expected, $s)) ", cls: c.dim },
        { text: "return ", cls: c.kw },
        { text: "true", cls: c.kw },
        { text: ";", cls: c.dim },
      ],
      [
        { text: "    }", cls: c.dim },
      ],
      [
        { text: "    return ", cls: c.kw },
        { text: "false", cls: c.kw },
        { text: ";", cls: c.dim },
      ],
      [
        { text: "}", cls: c.dim },
      ],
    ],
  },
];

function renderLine(segments: { text: string; cls?: string }[]) {
  return segments.map((seg, i) =>
    seg.cls ? (
      <span key={i} className={seg.cls}>{seg.text}</span>
    ) : (
      <span key={i}>{seg.text}</span>
    )
  );
}

export default function VerifyTabs() {
  const [activeTab, setActiveTab] = useState("nodejs");
  const [copied, setCopied] = useState(false);

  const activeSnippet = tabs.find((t) => t.id === activeTab)!;
  const plainText = activeSnippet.lines
    .map((segs) => segs.map((s) => s.text).join(""))
    .join("\n");

  const handleCopy = async () => {
    await navigator.clipboard.writeText(plainText);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="rounded-xl border border-surface-800 bg-surface-900 overflow-hidden">
      <div className="flex items-center border-b border-surface-800">
        <div className="flex-1 flex overflow-x-auto scrollbar-none">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => { setActiveTab(tab.id); setCopied(false); }}
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
        <code>
          {activeSnippet.lines.map((segs, i) => (
            <span key={i}>
              {segs.length === 0 ? "\n" : <>{renderLine(segs)}{"\n"}</>}
            </span>
          ))}
        </code>
      </pre>
    </div>
  );
}
