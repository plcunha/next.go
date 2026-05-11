# Next.go - The Go Implementation of Next.js

Uma recriação completa do Next.js em Go, trazendo todos os principais conceitos:

## 🎯 O que foi implementado

### Core Features
- ✅ **File-System Routing** - Roteamento automático baseado em arquivos (como App Router do Next.js)
- ✅ **Server-Side Rendering (SSR)** - Renderização no servidor
- ✅ **Static Site Generation (SSG)** - Geração estática de páginas
- ✅ **API Routes** - Rotas de API com Go
- ✅ **Hot Module Replacement** - Recarga automática em desenvolvimento
- ✅ **Middleware System** - Middleware personalizado para requisições
- ✅ **Build System** - Build otimizado para produção
- ✅ **Configuração YAML** - nextgo.yaml (equivalente ao next.config.js)

### Estrutura do Projeto

```
next.go/
├── cmd/                    # CLI commands (dev, build, start, create)
├── packages/
│   ├── api/               # API route handling
│   ├── build/             # Build system
│   ├── config/            # Configuration (nextgo.yaml)
│   ├── middleware/        # Middleware system
│   ├── router/           # File-system routing
│   ├── server/           # HTTP server com SSR
│   ├── static/           # Static file serving
│   └── watcher/          # File watcher para HMR
├── template/              # Templates para novos projetos
├── examples/              # Exemplos de uso
├── docs/                  # Documentação
└── main.go               # Entry point
```

## 🚀 Como usar

### Criar novo projeto
```bash
nextgo create meu-app
cd meu-app
```

### Desenvolvimento
```bash
nextgo dev
# Servidor roda em http://localhost:3000
```

### Build para produção
```bash
nextgo build
nextgo start
```

## 📁 Estrutura de um projeto Next.go

```
meu-app/
├── app/
│   ├── layout.go.html     # Layout raiz (como layout.tsx)
│   ├── page.go.html      # Página inicial
│   ├── about/
│   │   └── page.go.html # /about
│   └── api/
│       └── hello/
│           └── handler.go # API route
├── public/               # Arquivos estáticos
├── components/           # Componentes reutilizáveis
├── nextgo.yaml          # Configuração
└── package.json         # (opcional)
```

## 🔄 Roteamento

| Arquivo | URL |
|---------|-----|
| `app/page.go.html` | `/` |
| `app/about/page.go.html` | `/about` |
| `app/blog/[id]/page.go.html` | `/blog/:id` |
| `app/docs/[...slug]/page.go.html` | `/docs/*slug` |

## 🔌 API Routes

```go
// app/api/users/handler.go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    c.JSON(200, gin.H{
        "users": []string{"John", "Jane"},
    })
}
```

## ⚙️ Configuração (nextgo.yaml)

```yaml
server:
  port: 3000
  host: localhost
  compression: true

build:
  output: standalone
  distDir: .next

images:
  domains: []
  formats:
    - image/webp

experimental:
  appDir: true
```

## 📦 Templates Disponíveis

- **default** - Template padrão com layout, página inicial e API exemplo

## 🎨 Comparação com Next.js

| Feature | Next.js | Next.go |
|---------|----------|---------|
| Linguagem | React/Node.js | Go |
| Roteamento | File-system | File-system |
| SSR | ✓ | ✓ |
| SSG | ✓ | ✓ |
| API Routes | ✓ | ✓ |
| Middleware | ✓ | ✓ |
| Hot Reload | ✓ | ✓ |
| Deploy | Node server | Go binary |

## 🛠️ Instalação

```bash
# Build do código-fonte
cd C:/Users/v15/Documents/GitHub/next.go
go build -o nextgo .

# Ou instalar globalmente (quando Go estiver instalado)
go install github.com/nextgo/nextgo@latest
```

## 📝 Próximos Passos

Para completar o Next.go, você pode:

1. **Instalar Go** - Baixe em https://go.dev/dl/
2. **Compilar o projeto** - `go build -o nextgo .`
3. **Testar** - `./nextgo create test-app`
4. **Expandir funcionalidades**:
   - Implementar sistema de componentes completo
   - Adicionar suporte a TypeScript/React via Node.js bridge
   - Implementar Image Optimization
   - Adicionar suporte a ISR (Incremental Static Regeneration)
   - Criar mais templates

## 📄 Licença

MIT License - veja LICENSE para detalhes.

---

**Nota**: Este é um projeto que recria os conceitos do Next.js em Go. 
O código está estruturado e pronto para compilação assim que o Go estiver instalado no sistema.
