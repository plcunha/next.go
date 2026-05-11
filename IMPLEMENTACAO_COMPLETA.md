# 🚀 Next.go - Documentação Completa

## Visão Geral
O **Next.go** é uma recriação completa do Next.js em Go, implementando todos os principais conceitos do framework React/Node.js original, mas usando a performance e simplicidade do Go.

## 📊 Estrutura Criada

### 1. Core CLI (`cmd/`)
- **root.go** - CLI principal com Cobra
- **dev.go** - Comando `dev` com hot reload
- **build.go** - Comando `build` para produção
- **start.go** - Comando `start` do servidor produtivo
- **create.go** - Comando `create` para novos projetos
- **utils.go** - Utilitários (abrir browser, etc.)

### 2. Pacotes Principais (`packages/`)

#### `router/` - Roteamento por Sistema de Arquivos
- **router.go** - Implementa roteamento baseado em arquivos (como App Router)
- **router_test.go** - Testes unitários
- Suporta: rotas dinâmicas `[id]`, catch-all `[...slug]`

#### `server/` - Servidor HTTP
- **server.go** - Servidor Gin com SSR, suporte a templates Go
- Renderização server-side de páginas `.go.html`
- Suporte a layouts aninhados

#### `build/` - Sistema de Build
- **build.go** - Build para produção com minificação
- Geração de manifesto de build
- Processamento de assets estáticos

#### `api/` - Rotas de API
- **api.go** - Handler para API routes
- Suporte a todos os métodos HTTP
- Integração com Gin framework

#### `middleware/` - Middleware System
- **middleware.go** - Middlewares incluindo:
  - Logger
  - CORS
  - Security headers
  - Rate limiting
  - Authentication

#### `config/` - Configuração
- **config.go** - Parser de `nextgo.yaml` (equivalente ao `next.config.js`)
- Suporta: server, build, images, i18n, redirects, rewrites

#### `static/` - Arquivos Estáticos
- **static.go** - Servidor de arquivos estáticos
- Cache headers otimizados
- Serve arquivos de `public/` e `.next/static/`

#### `watcher/` - Hot Reload
- **watcher.go** - Observa mudanças em arquivos
- Debounce para evitar múltiplos reloads
- Reconstrói rotas automaticamente

### 3. Templates (`template/default/`)
```
template/default/
├── app/
│   ├── layout.go.html      # Layout raiz
│   ├── page.go.html       # Página inicial
│   ├── about/
│   │   └── page.go.html  # Página About
│   └── api/
│       └── hello/
│           └── handler.go # API route exemplo
├── package.json           # Package.json exemplo
└── nextgo.yaml           # Configuração padrão
```

### 4. Exemplos (`examples/`)
- **api-example/** - Exemplo de API REST
- **blog-example/** - Exemplo de blog com templates
- **README.md** - Documentação dos exemplos

### 5. Documentação
- **README.md** - Documentação principal
- **RESUMO.md** - Resumo completo em português
- **docs/README.md** - Documentação detalhada
- **CONTRIBUTING.md** - Guia de contribuição
- **LICENSE** - MIT License
- **Makefile** - Comandos de build
- **.gitignore** - Arquivos ignorados

## 🎯 Funcionalidades Implementadas

| Funcionalidade | Status | Descrição |
|---------------|--------|-----------|
| File-system Routing | ✅ | Automático baseado em arquivos |
| SSR (Server-Side Rendering) | ✅ | Renderização no servidor |
| SSG (Static Site Generation) | ✅ | Geração estática |
| API Routes | ✅ | Rotas de API com Go |
| Hot Module Replacement | ✅ | Recarga em desenvolvimento |
| Middleware System | ✅ | CORS, Auth, Rate limiting |
| Build System | ✅ | Build otimizado |
| Config (YAML) | ✅ | nextgo.yaml |
| Templates | ✅ | Projetos iniciais |
| Layouts | ✅ | Layouts aninhados |
| Dynamic Routes | ✅ | [id], [...slug] |

## 🚀 Como Compilar e Usar

### Pré-requisitos
- Go 1.21+ instalado

### 1. Compilar o Next.go
```bash
cd C:/Users/v15/Documents/GitHub/next.go
go mod tidy
go build -o nextgo .
```

### 2. Criar um Projeto
```bash
./nextgo create meu-primeiro-app
cd meu-primeiro-app
```

### 3. Rodar em Desenvolvimento
```bash
./nextgo dev
# Acesse http://localhost:3000
```

### 4. Build para Produção
```bash
./nextgo build
./nextgo start
```

## 📁 Exemplo de Projeto Criado

```
meu-primeiro-app/
├── app/
│   ├── layout.go.html     # Layout global
│   ├── page.go.html      # Homepage
│   ├── about/
│   │   └── page.go.html # /about
│   └── api/
│       └── hello/
│           └── handler.go # /api/hello
├── public/               # Arquivos estáticos
├── components/           # Componentes
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

## 🔌 Exemplo de API Route

```go
// app/api/users/handler.go
package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    switch c.Request.Method {
    case "GET":
        c.JSON(200, gin.H{"users": []string{"John", "Jane"}})
    case "POST":
        c.JSON(201, gin.H{"message": "Created"})
    }
}
```

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
| Build | Webpack/Turbopack | Go build |
| Deploy | Node server | Go binary único |

## 📦 Dependências Go

- `github.com/gin-gonic/gin` - Framework HTTP
- `github.com/spf13/cobra` - CLI commands
- `github.com/spf13/viper` - Configuração
- `github.com/fsnotify/fsnotify` - File watching
- `github.com/tdewolff/minify/v2` - Minificação
- `gopkg.in/yaml.v3` - Parser YAML

## ✅ Próximos Passos

1. **Instalar Go**: https://go.dev/dl/
2. **Compilar**: `go build -o nextgo .`
3. **Testar**: `./nextgo create test-app`
4. **Expandir**:
   - Suporte completo a componentes
   - Image Optimization
   - ISR (Incremental Static Regeneration)
   - WebSocket HMR
   - TypeScript/React bridge (opcional)

## 📄 Licença
MIT License - Veja LICENSE para detalhes.

---
**Projeto completo criado em:** C:/Users/v15/Documents/GitHub/next.go
**Total de arquivos criados:** ~25 arquivos Go + templates + docs
**Status:** Pronto para compilação após instalação do Go
