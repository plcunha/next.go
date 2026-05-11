# Next.go

A Go implementation of Next.js concepts — file-system routing, SSR, API routes, middleware, and production builds as a single binary.

## Features

- **File-system Routing** — Automatic routing based on directory structure (`[id]`, `[...slug]`)
- **Server-Side Rendering** — Go templates rendered on the server with layout wrapping
- **API Routes** — Write Go handlers, compiled and executed in dev and production
- **Middleware** — Logger, CORS, Security headers built-in
- **Hot Reload** — File watcher re-scans routes on change
- **Build System** — Minified HTML + compiled API binary (`app.exe`)
- **Single Binary Deploy** — `nextgo build` produces ~12MB binary, zero runtime deps

## Installation

```bash
go install github.com/nextgo/nextgo@latest
```

Or build from source:

```bash
git clone https://github.com/nextgo/nextgo.git
cd nextgo
go build -o nextgo .
```

## Quick Start

```bash
nextgo create my-app
cd my-app
nextgo dev                # http://localhost:3000
nextgo dev -p 8080        # custom port
```

## Build & Production

```bash
nextgo build              # compiles pages + API handlers → .next/
nextgo start              # serves from .next/
nextgo start -p 8080      # custom port
```

## Project Structure

```
my-app/
├── app/
│   ├── layout.go.html      # Root layout
│   ├── page.go.html        # → /
│   ├── about/
│   │   └── page.go.html    # → /about
│   └── api/
│       └── hello/
│           └── handler.go  # → /api/hello
├── public/                 # Static assets
├── nextgo.yaml             # Configuration
└── package.json            # npm scripts (optional)
```

## Routing

| File | URL |
|---|---|
| `app/page.go.html` | `/` |
| `app/about/page.go.html` | `/about` |
| `app/blog/[id]/page.go.html` | `/blog/:id` |
| `app/docs/[...slug]/page.go.html` | `/docs/*slug` |
| `app/api/hello/handler.go` | `/api/hello` |

## API Routes

```go
// app/api/hello/handler.go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Hello from Next.go!",
        "method":  c.Request.Method,
    })
}
```

In dev mode, handlers are compiled into a subprocess and executed in real time.
In production, handlers are compiled into the `app.exe` binary.

## Templates

`.go.html` files are Go templates with automatic layout wrapping:

```html
<!-- app/layout.go.html -->
<html>
<head><title>{{.title}}</title></head>
<body>
    <nav><a href="/">Home</a> | <a href="/about">About</a></nav>
    <main>{{.children}}</main>
</body>
</html>
```

Available template variables: `{{.title}}`, `{{.path}}`, `{{.children}}` (layout only).

## Configuration

```yaml
# nextgo.yaml
server:
  port: 3000
  host: localhost

build:
  output: standalone
  distDir: .next
```

## Middleware

Active in dev and production:

| Middleware | Headers |
|---|---|
| Recovery | Panic recovery |
| Logger | `[NEXT-GO] GET / 200 1.1ms` |
| CORS | `Access-Control-Allow-Origin: *` |
| Security | `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `X-XSS-Protection`, `Referrer-Policy` |

## Feature Status

| Feature | Status |
|---|---|
| File-system Routing | ✅ |
| SSR (dev) | ✅ |
| SSG (build) | ✅ |
| API Routes (dev) | ✅ |
| API Routes (build) | ✅ |
| Middleware | ✅ |
| Hot Reload | ✅ |
| Build → Binary | ✅ |
| `go install` | ✅ |

## Documentation

See [docs/README.md](docs/README.md) for the complete guide.

## License

MIT
