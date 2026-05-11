# Next.go Documentation

Welcome to Next.go - A Go implementation of Next.js framework concepts.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Project Structure](#project-structure)
3. [Routing](#routing)
4. [Pages and Layouts](#pages-and-layouts)
5. [API Routes](#api-routes)
6. [Configuration](#configuration)
7. [CLI Commands](#cli-commands)
8. [Build and Deployment](#build-and-deployment)

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

```bash
go install github.com/nextgo/nextgo@latest
```

Or build from source:

```bash
git clone https://github.com/nextgo/nextgo.git
cd nextgo
go build -o nextgo .
```

### Create Your First App

```bash
nextgo create my-app
cd my-app
nextgo dev
```

Visit `http://localhost:3000` to see your app.

## Project Structure

```
my-app/
├── app/                    # App directory (like Next.js App Router)
│   ├── layout.go.html     # Root layout
│   ├── page.go.html      # Home page
│   ├── about/
│   │   └── page.go.html  # About page
│   └── api/              # API routes
│       └── hello/
│           └── handler.go
├── public/               # Static assets
├── components/           # Reusable components
├── lib/                  # Utility functions
├── nextgo.yaml          # Configuration
└── package.json         # (optional) for Node.js tooling
```

## Routing

Next.go uses file-system routing similar to Next.js App Router.

### Basic Routes

| File | URL Path |
|------|----------|
| `app/page.go.html` | `/` |
| `app/about/page.go.html` | `/about` |
| `app/blog/page.go.html` | `/blog` |

### Dynamic Routes

| File | URL Path |
|------|----------|
| `app/blog/[id]/page.go.html` | `/blog/:id` |
| `app/docs/[...slug]/page.go.html` | `/docs/*slug` |

## Pages and Layouts

### Pages

Pages are `.go.html` files that contain HTML with Go template syntax:

```html
<!-- app/about/page.go.html -->
<div class="about">
    <h1>About Us</h1>
    <p>Welcome to our website!</p>
</div>

<style>
.about h1 { color: #000; }
</style>
```

### Layouts

Layouts wrap pages and persist across navigations:

```html
<!-- app/layout.go.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
</head>
<body>
    <nav><!-- navigation --></nav>
    <main>{{.children}}</main>
    <footer><!-- footer --></footer>
</body>
</html>
```

## API Routes

Create API routes in `app/api` directory:

```go
// app/api/users/handler.go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    switch c.Request.Method {
    case "GET":
        c.JSON(200, gin.H{"users": []string{"John", "Jane"}})
    case "POST":
        c.JSON(201, gin.H{"message": "User created"})
    }
}
```

## Configuration

Configure your app with `nextgo.yaml`:

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

## CLI Commands

### `nextgo dev`

Starts the development server with hot reload.

```bash
nextgo dev
# Options:
#   -p, --port <port>    Port to run the server (default: 3000)
```

### `nextgo build`

Builds the app for production.

```bash
nextgo build
```

Output will be in `.next/` directory.

### `nextgo start`

Starts the production server.

```bash
nextgo start
```

### `nextgo create <project-name>`

Creates a new Next.go project.

```bash
nextgo create my-app
```

## Build and Deployment

### Build for Production

```bash
nextgo build
```

### Start Production Server

```bash
nextgo start
# Or with custom port
PORT=8080 nextgo start
```

### Deploy

Since Next.go compiles to a single binary, deployment is simple:

1. Build for your target platform
2. Copy the binary and `.next` directory
3. Run the binary

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o nextgo .

# Copy to server
scp nextgo user@server:/app/
scp -r .next user@server:/app/

# On server
cd /app && ./nextgo start
```

## Comparison with Next.js

| Feature | Next.js | Next.go |
|---------|---------|---------|
| Language | JavaScript/TypeScript + React | Go + HTML Templates |
| Routing | File-system | File-system |
| SSR | ✓ | ✓ |
| SSG | ✓ | ✓ |
| API Routes | ✓ | ✓ |
| Middleware | ✓ | ✓ |
| Hot Reload | ✓ | ✓ |
| Build Output | Node.js server | Go binary |

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.
