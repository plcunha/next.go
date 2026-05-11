# Next.go - Guia Completo

Next.go é uma implementação em Go dos conceitos de **file-system routing** e **server-side rendering** do Next.js. Aplicações web com o paradigma do Next.js — Go + templates HTML, sem JavaScript no frontend.

---

## Instalação

### Pré-requisitos

- Go 1.21+ ([go.dev/dl](https://go.dev/dl/))
- `$GOPATH/bin` no PATH:

```powershell
# Windows - permanente
[Environment]::SetEnvironmentVariable("PATH", "$env:PATH;$env:USERPROFILE\go\bin", "User")
```

```bash
# Linux/macOS
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
```

### Instalar

```bash
# Da rede (repo público):
go install github.com/nextgo/nextgo@latest

# Ou do source local:
cd next.go
go install .
```

### Verificar

```bash
nextgo --help
```

---

## Quick Start

```bash
nextgo create my-app
cd my-app
nextgo dev              # http://localhost:3000
nextgo dev -p 8080      # porta customizada
```

### Estrutura do projeto

```
my-app/
├── app/
│   ├── layout.go.html      # Layout global
│   ├── page.go.html        # → /
│   ├── about/
│   │   └── page.go.html    # → /about
│   └── api/
│       └── hello/
│           └── handler.go  # → /api/hello
├── public/                 # Arquivos estáticos
├── nextgo.yaml             # Configuração
└── package.json            # Scripts (npm scripts opcionais)
```

---

## Rotas

### File-system routing

| Arquivo | Rota | Tipo |
|---|---|---|
| `app/page.go.html` | `/` | Page SSR |
| `app/about/page.go.html` | `/about` | Page SSR |
| `app/blog/[id]/page.go.html` | `/blog/:id` | Dinâmica |
| `app/docs/[...slug]/page.go.html` | `/docs/*slug` | Catch-all |
| `app/api/hello/handler.go` | `/api/hello` | API |

### Ignorados

`layout.go.html`, `loading.go.html`, `error.go.html`, `not-found.go.html`, arquivos com `_` ou `.`

---

## Páginas e Templates

### Page

```html
<div><h1>Minha Página</h1><p>Conteúdo SSR.</p></div>
```

### Layout

```html
<html>
<head><title>{{.title}}</title></head>
<body>
    <nav><a href="/">Home</a> | <a href="/about">About</a></nav>
    <main>{{.children}}</main>
</body>
</html>
```

O layout envolve páginas do mesmo nível ou abaixo automaticamente.

---

## API Routes

`app/api/<rota>/handler.go`:

```go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    switch c.Request.Method {
    case "GET":
        c.JSON(200, gin.H{"users": []string{"Alice", "Bob"}})
    case "POST":
        var body map[string]interface{}
        c.ShouldBindJSON(&body)
        c.JSON(201, gin.H{"created": body})
    default:
        c.JSON(405, gin.H{"error": "Method not allowed"})
    }
}
```

**No dev mode:** handlers são compilados automaticamente num subprocesso e executados em tempo real.
**No build:** handlers são compilados no binário `app.exe`.

---

## Build

```bash
nextgo build
```

O build:
1. Minifica páginas HTML → `.next/server/app/`
2. Copia handlers API → `.next/build/handlers/`
3. Gera `main.go` + `go.mod`
4. Compila `app.exe` → `.next/server/app.exe` (~12MB)
5. Copia assets estáticos → `.next/static/`
6. Gera `build-manifest.json`

### Executar produção

```bash
nextgo start              # :3000
nextgo start -p 8080      # :8080

# Ou direto o binário compilado:
.next/server/app.exe      # :3000 (ou $PORT)
```

---

## Middleware

Ativado automaticamente em dev e produção:

| Middleware | Função |
|---|---|
| `gin.Recovery()` | Recuperação de panics |
| `middleware.Logger()` | Log de requests `[NEXT-GO] GET / 200 1.1ms` |
| `middleware.CORS()` | Headers CORS completos |
| `middleware.Security()` | `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Referrer-Policy` |

Headers verificados:
```
Access-Control-Allow-Origin: *
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
```

---

## Configuração

`nextgo.yaml`:

```yaml
server:
  port: 3000
  host: localhost

build:
  output: standalone
  distDir: .next
```

---

## CLI Reference

| Comando | Descrição |
|---|---|
| `nextgo create <nome>` | Cria novo projeto |
| `nextgo dev [-p port]` | Dev server com hot reload |
| `nextgo build` | Build produção (páginas + API compilada) |
| `nextgo start [-p port]` | Server produção |
| `nextgo --help` | Ajuda |

---

## Status das Features

| Feature | Status | Detalhes |
|---|---|---|
| File-system Routing | ✅ | `[id]`, `[...slug]`, layouts |
| SSR (dev) | ✅ | Templates Go com layout wrapping |
| SSG (build) | ✅ | HTML minificado + manifest |
| API Routes (build) | ✅ | Compiladas em `app.exe` |
| API Routes (dev) | ✅ | Compilados em subprocesso, executados em tempo real |
| Middleware | ✅ | Logger + CORS + Security |
| Hot Reload | ✅ | Watcher + re-scan |
| Build → Binary | ✅ | `app.exe` ~12MB |
| `go install` | ✅ | Global via `$GOPATH/bin` |
| Templates no build | ⚠️ | HTML concatenado; `{{.var}}` não executado |

---

## Quando usar ✅

- Blogs, portfolios, landing pages
- APIs com páginas de admin
- Microserviços com UI leve
- Performance/memória crítica
- Times Go que querem SSR sem Node.js
- Deploy com 1 binário

## Quando não usar ❌

- Apps web interativos (precisa React/SPA)
- UIs complexas com estado
- E-commerce com carrinho
- PWA mobile-first

---

## Benchmark

| Métrica | Next.js | Next.go |
|---|---|---|
| Startup | 2-5s | ~50ms |
| Memória idle | 150-300MB | 10-20MB |
| Binário | node_modules ~300MB | ~20MB |
| Deploy | npm install + build | 1 binário |

---

## Roadmap

### ✅ Implementado
- File-system routing (`[id]`, `[...slug]`)
- SSR com layout wrapping
- Build com minificação HTML
- API routes compiladas em `app.exe` (build)
- API routes compiladas em subprocesso (dev)
- Middleware chain (Logger, CORS, Security)
- Build manifest
- Hot reload
- `go install` global

### ⏳ Pendente
2. Template execution no build
3. Hybrid mode (JS frontend)
4. htmx/Alpine.js support
5. Partials/slots/templates
6. Image optimization
7. ISR
8. Database integration
9. Authentication (OAuth, JWT)
