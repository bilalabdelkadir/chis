import { useState } from "react";

interface Tab {
  id: string;
  label: string;
  lines: { text: string; cls?: string }[][];
}

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

const tabs: Tab[] = [
  {
    id: "typescript",
    label: "TypeScript",
    lines: [
      [
        { text: "const ", cls: c.kw },
        { text: "res ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "await ", cls: c.kw },
        { text: "fetch", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ", {", cls: c.dim },
      ],
      [
        { text: "  method", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: '"POST"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  headers", cls: c.key },
        { text: ": {", cls: c.dim },
      ],
      [
        { text: '    "Authorization"', cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "`Bearer ${", cls: c.str },
        { text: "process.env.", cls: c.plain },
        { text: "CHIS_API_KEY", cls: c.key },
        { text: "}`", cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '    "Content-Type"', cls: c.key },
        { text: ":  ", cls: c.dim },
        { text: '"application/json"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  },", cls: c.dim },
      ],
      [
        { text: "  body", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "JSON.", cls: c.plain },
        { text: "stringify", cls: c.fn },
        { text: "({", cls: c.dim },
      ],
      [
        { text: "    url", cls: c.key },
        { text: ":     ", cls: c.dim },
        { text: '"https://your-app.com/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    event", cls: c.key },
        { text: ":   ", cls: c.dim },
        { text: '"payment.completed"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    payload", cls: c.key },
        { text: ": { ", cls: c.dim },
        { text: "amount", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "4999", cls: c.num },
        { text: ", ", cls: c.dim },
        { text: "currency", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: '"usd"', cls: c.str },
        { text: " },", cls: c.dim },
      ],
      [
        { text: "  }),", cls: c.dim },
      ],
      [
        { text: "});", cls: c.dim },
      ],
      [],
      [
        { text: "const ", cls: c.kw },
        { text: "data", cls: c.plain },
        { text: ": ", cls: c.dim },
        { text: "Delivery", cls: c.fn },
        { text: " = ", cls: c.dim },
        { text: "await ", cls: c.kw },
        { text: "res.", cls: c.plain },
        { text: "json", cls: c.fn },
        { text: "();", cls: c.dim },
      ],
    ],
  },
  {
    id: "nodejs",
    label: "Node.js",
    lines: [
      [
        { text: "const ", cls: c.kw },
        { text: "res ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "await ", cls: c.kw },
        { text: "fetch", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ", {", cls: c.dim },
      ],
      [
        { text: "  method", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: '"POST"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  headers", cls: c.key },
        { text: ": {", cls: c.dim },
      ],
      [
        { text: '    "Authorization"', cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "`Bearer ${", cls: c.str },
        { text: "process.env.", cls: c.plain },
        { text: "CHIS_API_KEY", cls: c.key },
        { text: "}`", cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '    "Content-Type"', cls: c.key },
        { text: ":  ", cls: c.dim },
        { text: '"application/json"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  },", cls: c.dim },
      ],
      [
        { text: "  body", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "JSON.", cls: c.plain },
        { text: "stringify", cls: c.fn },
        { text: "({", cls: c.dim },
      ],
      [
        { text: "    url", cls: c.key },
        { text: ":     ", cls: c.dim },
        { text: '"https://your-app.com/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    event", cls: c.key },
        { text: ":   ", cls: c.dim },
        { text: '"payment.completed"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    payload", cls: c.key },
        { text: ": { ", cls: c.dim },
        { text: "amount", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "4999", cls: c.num },
        { text: " }", cls: c.dim },
      ],
      [
        { text: "  }),", cls: c.dim },
      ],
      [
        { text: "});", cls: c.dim },
      ],
      [],
      [
        { text: "const ", cls: c.kw },
        { text: "data ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "await ", cls: c.kw },
        { text: "res.", cls: c.plain },
        { text: "json", cls: c.fn },
        { text: "();", cls: c.dim },
      ],
    ],
  },
  {
    id: "python",
    label: "Python",
    lines: [
      [
        { text: "import ", cls: c.kw },
        { text: "requests", cls: c.plain },
      ],
      [
        { text: "import ", cls: c.kw },
        { text: "os", cls: c.plain },
      ],
      [],
      [
        { text: "response ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "requests.", cls: c.plain },
        { text: "post", cls: c.fn },
        { text: "(", cls: c.dim },
      ],
      [
        { text: '    "https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    headers", cls: c.key },
        { text: "={", cls: c.dim },
      ],
      [
        { text: '        "Authorization"', cls: c.key },
        { text: ": ", cls: c.dim },
        { text: 'f"Bearer {', cls: c.str },
        { text: "os.environ[", cls: c.plain },
        { text: "'CHIS_API_KEY'", cls: c.str },
        { text: "]", cls: c.plain },
        { text: '}"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '        "Content-Type"', cls: c.key },
        { text: ":  ", cls: c.dim },
        { text: '"application/json"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "    },", cls: c.dim },
      ],
      [
        { text: "    json", cls: c.key },
        { text: "={", cls: c.dim },
      ],
      [
        { text: '        "url"', cls: c.key },
        { text: ":     ", cls: c.dim },
        { text: '"https://your-app.com/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '        "event"', cls: c.key },
        { text: ":   ", cls: c.dim },
        { text: '"payment.completed"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '        "payload"', cls: c.key },
        { text: ": {", cls: c.dim },
        { text: '"amount"', cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "4999", cls: c.num },
        { text: "},", cls: c.dim },
      ],
      [
        { text: "    },", cls: c.dim },
      ],
      [
        { text: ")", cls: c.dim },
      ],
      [],
      [
        { text: "print", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "response.", cls: c.plain },
        { text: "json", cls: c.fn },
        { text: "())", cls: c.dim },
      ],
    ],
  },
  {
    id: "go",
    label: "Go",
    lines: [
      [
        { text: "body ", cls: c.plain },
        { text: ":= ", cls: c.dim },
        { text: "strings.", cls: c.plain },
        { text: "NewReader", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: "`{", cls: c.str },
      ],
      [
        { text: '  "url":     "https://your-app.com/webhooks",', cls: c.str },
      ],
      [
        { text: '  "event":   "payment.completed",', cls: c.str },
      ],
      [
        { text: '  "payload": {"amount": 4999}', cls: c.str },
      ],
      [
        { text: "}`", cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [],
      [
        { text: "req, _ ", cls: c.plain },
        { text: ":= ", cls: c.dim },
        { text: "http.", cls: c.plain },
        { text: "NewRequest", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"POST"', cls: c.str },
        { text: ", ", cls: c.dim },
        { text: '"https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ", body)", cls: c.dim },
      ],
      [
        { text: "req.", cls: c.plain },
        { text: "Header.", cls: c.plain },
        { text: "Set", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"Authorization"', cls: c.str },
        { text: ", ", cls: c.dim },
        { text: '"Bearer "', cls: c.str },
        { text: " + ", cls: c.dim },
        { text: "os.", cls: c.plain },
        { text: "Getenv", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"CHIS_API_KEY"', cls: c.str },
        { text: "))", cls: c.dim },
      ],
      [
        { text: "req.", cls: c.plain },
        { text: "Header.", cls: c.plain },
        { text: "Set", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"Content-Type"', cls: c.str },
        { text: ", ", cls: c.dim },
        { text: '"application/json"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [],
      [
        { text: "resp, _ ", cls: c.plain },
        { text: ":= ", cls: c.dim },
        { text: "http.", cls: c.plain },
        { text: "DefaultClient.", cls: c.plain },
        { text: "Do", cls: c.fn },
        { text: "(req)", cls: c.dim },
      ],
      [
        { text: "defer ", cls: c.kw },
        { text: "resp.", cls: c.plain },
        { text: "Body.", cls: c.plain },
        { text: "Close", cls: c.fn },
        { text: "()", cls: c.dim },
      ],
    ],
  },
  {
    id: "ruby",
    label: "Ruby",
    lines: [
      [
        { text: "require ", cls: c.kw },
        { text: '"net/http"', cls: c.str },
      ],
      [
        { text: "require ", cls: c.kw },
        { text: '"json"', cls: c.str },
      ],
      [],
      [
        { text: "uri ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "URI", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ")", cls: c.dim },
      ],
      [
        { text: "req ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "Net::HTTP::Post.", cls: c.plain },
        { text: "new", cls: c.fn },
        { text: "(uri)", cls: c.dim },
      ],
      [],
      [
        { text: 'req["Authorization"] ', cls: c.key },
        { text: "= ", cls: c.dim },
        { text: '"Bearer #{', cls: c.str },
        { text: "ENV[", cls: c.plain },
        { text: "'CHIS_API_KEY'", cls: c.str },
        { text: "]", cls: c.plain },
        { text: '}"', cls: c.str },
      ],
      [
        { text: 'req["Content-Type"]  ', cls: c.key },
        { text: "= ", cls: c.dim },
        { text: '"application/json"', cls: c.str },
      ],
      [
        { text: "req.", cls: c.plain },
        { text: "body ", cls: c.key },
        { text: "= {", cls: c.dim },
      ],
      [
        { text: "  url", cls: c.key },
        { text: ":     ", cls: c.dim },
        { text: '"https://your-app.com/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  event", cls: c.key },
        { text: ":   ", cls: c.dim },
        { text: '"payment.completed"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: "  payload", cls: c.key },
        { text: ": { ", cls: c.dim },
        { text: "amount", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "4999", cls: c.num },
        { text: " }", cls: c.dim },
      ],
      [
        { text: "}.", cls: c.dim },
        { text: "to_json", cls: c.fn },
      ],
      [],
      [
        { text: "res ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "Net::HTTP.", cls: c.plain },
        { text: "start", cls: c.fn },
        { text: "(uri.host, uri.port, ", cls: c.dim },
        { text: "use_ssl", cls: c.key },
        { text: ": ", cls: c.dim },
        { text: "true", cls: c.kw },
        { text: ") { |h| h.", cls: c.dim },
        { text: "request", cls: c.fn },
        { text: "(req) }", cls: c.dim },
      ],
      [
        { text: "puts ", cls: c.fn },
        { text: "res.", cls: c.plain },
        { text: "body", cls: c.key },
      ],
    ],
  },
  {
    id: "php",
    label: "PHP",
    lines: [
      [
        { text: "$response ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "Http", cls: c.fn },
        { text: "::", cls: c.dim },
        { text: "withHeaders", cls: c.fn },
        { text: "([", cls: c.dim },
      ],
      [
        { text: '    "Authorization" ', cls: c.key },
        { text: "=> ", cls: c.dim },
        { text: '"Bearer "', cls: c.str },
        { text: " . ", cls: c.dim },
        { text: "env", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"CHIS_API_KEY"', cls: c.str },
        { text: "),", cls: c.dim },
      ],
      [
        { text: "])", cls: c.dim },
        { text: "->", cls: c.dim },
        { text: "post", cls: c.fn },
        { text: "(", cls: c.dim },
        { text: '"https://api.trychis.com/v1/webhooks"', cls: c.str },
        { text: ", [", cls: c.dim },
      ],
      [
        { text: '    "url"     ', cls: c.key },
        { text: "=> ", cls: c.dim },
        { text: '"https://your-app.com/webhooks"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '    "event"   ', cls: c.key },
        { text: "=> ", cls: c.dim },
        { text: '"payment.completed"', cls: c.str },
        { text: ",", cls: c.dim },
      ],
      [
        { text: '    "payload" ', cls: c.key },
        { text: "=> [", cls: c.dim },
        { text: '"amount"', cls: c.key },
        { text: " => ", cls: c.dim },
        { text: "4999", cls: c.num },
        { text: "],", cls: c.dim },
      ],
      [
        { text: "]);", cls: c.dim },
      ],
      [],
      [
        { text: "$data ", cls: c.plain },
        { text: "= ", cls: c.dim },
        { text: "$response", cls: c.plain },
        { text: "->", cls: c.dim },
        { text: "json", cls: c.fn },
        { text: "();", cls: c.dim },
      ],
    ],
  },
  {
    id: "curl",
    label: "cURL",
    lines: [
      [
        { text: "curl", cls: c.fn },
        { text: " -X POST ", cls: c.dim },
        { text: "https://api.trychis.com/v1/webhooks", cls: c.plain },
        { text: " \\", cls: c.dim },
      ],
      [
        { text: "  -H ", cls: c.dim },
        { text: '"Authorization: Bearer $CHIS_API_KEY"', cls: c.str },
        { text: " \\", cls: c.dim },
      ],
      [
        { text: "  -H ", cls: c.dim },
        { text: '"Content-Type: application/json"', cls: c.str },
        { text: " \\", cls: c.dim },
      ],
      [
        { text: "  -d ", cls: c.dim },
        { text: "'{", cls: c.str },
      ],
      [
        { text: '    "url":     "https://your-app.com/webhooks",', cls: c.str },
      ],
      [
        { text: '    "event":   "payment.completed",', cls: c.str },
      ],
      [
        { text: '    "payload": { "amount": 4999 }', cls: c.str },
      ],
      [
        { text: "  }'", cls: c.str },
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

export default function CodeTabs() {
  const [activeTab, setActiveTab] = useState("typescript");
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
