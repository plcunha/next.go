# Next.go

A Go implementation of Next.js concepts — file-system routing, SSR, API routes, middleware, and production builds as a single binary.

## Features

- **File-system Routing** — Automatic routing based on directory structure (`[id]`, `[...slug]`)
- **Server-Side Rendering** — Go templates rendered on the server with layout wrapping and optional data fetching (`getServerSideProps`)
- **API Routes** — Write Go handlers, compiled and executed in dev (subprocess) and production (single binary)
- **Middleware** — Logger, CORS, Security headers built-in
- **Hot Reload** — File watcher re-scans routes on change
- **Build System** — Minified SSG pages + compiled API binary (~11MB)
- **Single Binary Deploy** — `nextgo build` produces standalone binary, zero runtime deps
- **Embedded Templates** — Default project templates baked into the binary via `go:embed`

## Installation

```bash
go install github.com/plcunha/next.go@latest
```

Or build from source:

```bash
git clone https://github.com/plcunha/next.go.git
cd next.go
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
                          # produces .next/server/app (standalone binary)

# Run the compiled binary directly:
.next/server/app          # http://localhost:3000
PORT=8080 .next/server/app  # custom port

# Or use the dev server in production mode:
nextgo start -p 8080
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
In production (`nextgo build`), handlers are compiled into a single standalone binary.

## Templates & SSR

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

Available template variables: `{{.title}}`, `{{.path}}`, `{{.route}}`, `{{.children}}` (layout only).

### getServerSideProps (data fetching)

Pages can fetch data at request time by registering a `Props` function:

```go
// app/posts/[id]/page.go (alongside page.go.html)
package posts

import "github.com/gin-gonic/gin"

func Props(c *gin.Context) gin.H {
    id := c.Param("id")
    post := fetchPost(id) // your data fetching logic
    return gin.H{"post": post, "title": post.Title}
}
```

Register in your `main.go` or init:

```go
s := server.New(".")
s.RegisterProps("/posts/:id", posts.Props)
```

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
| SSR with layout wrapping | ✅ |
| SSG (build → minified HTML) | ✅ |
| API Routes (dev, subprocess) | ✅ |
| API Routes (build, single binary) | ✅ |
| getServerSideProps | ✅ |
| Middleware | ✅ |
| Hot Reload | ✅ |
| Build → Binary (~11MB) | ✅ |
| Embedded templates (go:embed) | ✅ |
| `go test` suite (7 tests) | ✅ |

## Documentation

See [docs/README.md](docs/README.md) for the complete guide.

## License

MIT
