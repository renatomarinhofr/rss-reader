# RSS Reader Monolito (Go + React)

Aplicação monolítica que combina Go (backend) e React (frontend) para consumir, normalizar e exibir feeds RSS brasileiros de forma prática. O backend segue princípios de Clean Architecture, faz cache das consultas em PostgreSQL e expõe uma API simples; o frontend, construído com Vite + React, oferece catálogo de fontes, tema dark responsivo e leitura amigável do conteúdo original (incluindo imagens).

## Tecnologias e decisões

- **Go 1.23** com `net/http` + Clean Architecture: casos de uso isolados (`fetchfeed`, `listfeeds`, `clearfeeds`), interfaces para facilitar testes e injeção de dependência.
- **Gofeed** para parsear RSS/Atom, lidando com diferentes formatos de feeds brasileiros.
- **PostgreSQL + pgx** para armazenar snapshot dos feeds (cache) e histórico recente.
- **React 18 + Vite + TypeScript** para uma UI rápida, com hooks customizados (`useFeed`, `useRecentFeeds`) e catálogo pré-curado de fontes nacionais.
- **DOMPurify** sanitiza o HTML retornado pelos feeds, permitindo renderizar imagens, links e formatação com segurança.
- **Docker multi-stage** para gerar imagem mínima (distroless) e `docker compose` orquestrando app + banco.

## Estrutura

```
.
├── cmd/server         # Ponto de entrada do binário HTTP
├── internal
│   ├── domain/feed    # Entidades de domínio
│   ├── infra          # Implementações de infraestrutura (HTTP client, Postgres, ...)
│   ├── interface/http # Handlers e servidor HTTP
│   ├── repository     # Contratos de acesso a dados
│   └── usecase        # Casos de uso (aplicação)
└── web                # Frontend React com Vite
```

## Executando localmente

### Dependências

- Go 1.23
- Node.js 20+
- PostgreSQL 16 (ou ajustar `DATABASE_URL`)

### Variáveis de ambiente

| Nome | Descrição | Valor padrão |
|------|-----------|--------------|
| `PORT` | Porta HTTP do servidor Go | `8080` |
| `DATABASE_URL` | String de conexão PostgreSQL | `postgres://postgres:postgres@localhost:5432/rssreader?sslmode=disable` |

### Backend

```bash
go run ./cmd/server
```

A API ficará disponível em `http://localhost:8080` com os endpoints:

- `GET /api/feed?url=https://...` — busca o feed (usa cache se o download falhar) e persiste a última versão.
- `GET /api/feeds/recent` — lista os últimos feeds consultados armazenados no banco.
- `DELETE /api/feeds/recent` — limpa o histórico armazenado.
- `GET /healthz` — verificação de saúde.

O caso de uso de busca utiliza a biblioteca [`mmcdole/gofeed`](https://github.com/mmcdole/gofeed) para normalizar RSS/Atom.

### Testes

```bash
go test ./...
```

### Frontend (React)

No diretório `web`:

```bash
npm install
npm run dev
```

O Vite sobe em `http://localhost:5173` e utiliza proxy para o backend (porta 8080). Para build de produção:

```bash
npm run build
```

Os arquivos gerados em `web/dist` são servidos automaticamente pelo backend quando presentes, permitindo executar apenas o binário Go em produção.

## Ambiente com Docker

É possível subir o sistema completo via Docker com:

```bash
docker compose up --build
```

Os serviços expostos:

- `rssreader-app`: binário Go servindo API e frontend estático em `http://localhost:8080`.
- `rssreader-db`: instância PostgreSQL 16 com volume persistente (`postgres_data`).

Para parar e remover os containers:

```bash
docker compose down
```

## Fluxo de como tudo se conecta

1. O usuário informa a URL do feed no frontend React.
2. O frontend chama `GET /api/feed?url=...` no backend.
3. O caso de uso `fetchfeed` baixa o XML via `FeedFetcher`, normaliza com o **gofeed** e registra o snapshot (com hora) no PostgreSQL.
4. Os feeds consultados são armazenados e podem ser listados em `/api/feeds/recent`; em caso de erro de rede, o cache é usado como fallback.
5. O handler expõe DTOs JSON que a UI consome para montar cards com HTML sanitizado.

## Roadmap sugerido

- Autenticação + favoritos por usuário.
- Scheduler para atualização periódica de feeds e alertas.
- Observabilidade (Prometheus/OpenTelemetry) e métricas por fonte.
- Internacionalização e exportação (por exemplo, salvar em formato EPUB/PDF).
