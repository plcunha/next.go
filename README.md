# Next.go

A Go implementation of Next.js framework concepts - file-system routing, SSR, SSG, API routes, and more.

## Features

- 🚀 **File-system Routing** - Automatic routing based on file structure (like Next.js App Router)
- ⚡ **Server-Side Rendering (SSR)** - Render pages on the server
- 📄 **Static Site Generation (SSG)** - Pre-render pages at build time  
- 🔌 **API Routes** - Build API endpoints with Go
- 🔥 **Hot Module Replacement** - Fast refresh in development
- 🎨 **Middleware Support** - Custom middleware for requests
- 📦 **Build System** - Optimized production builds
- 🔒 **Security Built-in** - CORS, rate limiting, security headers

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

### Create a new project

```bash
nextgo create my-app
cd my-app
```

### Start development server

```bash
nextgo dev
```

The server will start at `http://localhost:3000`

### Build for production

```bash
nextgo build
```

### Start production server

```bash
nextgo start
```

## Project Structure

```
my-app/
├── app/
│   ├── layout.go.html    # Root layout (like app/layout.tsx)
│   ├── page.go.html     # Home page (like app/page.tsx)
│   ├── about/
│   │   └── page.go.html # About page
│   └── api/
│       └── hello/
│           └── handler.go # API route
├── public/              # Static assets
├── components/          # Reusable components
└── nextgo.yaml         # Configuration (like next.config.js)
```

## File-System Routing

Next.go uses file-system routing similar to Next.js:

| File | URL |
|------|-----|
| `app/page.go.html` | `/` |
| `app/about/page.go.html` | `/about` |
| `app/blog/[id]/page.go.html` | `/blog/:id` |
| `app/blog/[...slug]/page.go.html` | `/blog/*slug` |

## API Routes

Create API routes in the `app/api` directory:

```go
// app/api/hello/handler.go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Hello from Next.go!",
    })
}
```

## Configuration

Create a `nextgo.yaml` file for configuration:

```yaml
server:
  port: 3000
  host: localhost
  compression: true

build:
  output: standalone
  distDir: .next

images:
  domains:
    - example.com
  formats:
    - image/webp

experimental:
  appDir: true
```

## Templates

Next.go uses `.go.html` files for pages, which are Go templates with HTML:

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
</head>
<body>
    <h1>{{.title}}</h1>
    <p>{{.content}}</p>
</body>
</html>
```

## Comparison with Next.js

| Feature | Next.js | Next.go |
|---------|---------|---------|
| Language | React/Node.js | Go |
| Routing | File-system | File-system |
| SSR | ✓ | ✓ |
| SSG | ✓ | ✓ |
| API Routes | ✓ | ✓ |
| Middleware | ✓ | ✓ |
| Hot Reload | ✓ | ✓ |

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
