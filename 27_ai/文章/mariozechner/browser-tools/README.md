# Browser Tools

Minimal CDP tools for collaborative site exploration.

## Start Chrome

```bash
./start.js # Fresh profile
./start.js --profile # Copy your profile (cookies, logins)
```

Start Chrome on `:9222` with remote debugging.

## Navigate

```bash
./nav.js https://example.com
./nav.js https://example.com --new
```

Navigate current tab or open new tab.

## Evaluate JavaScript

```bash
./eval.js 'document.title'
./eval.js 'document.querySelectorAll("a").length'
```

Execute JavaScript in active tab (async context).

## Screenshot

```bash
./screenshot.js
```

Screenshot current viewport, returns temp file path.

## Pick Elements

```bash
./pick.js "Click the submit button"
```

Interactive element picker. Click to select, Cmd/Ctrl+Click for multi-select, Enter to finish.
