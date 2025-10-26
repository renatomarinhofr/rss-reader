# syntax=docker/dockerfile:1

FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web .
RUN npm run build

FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o rss-reader ./cmd/server

FROM gcr.io/distroless/base-debian12 AS runtime
WORKDIR /app
COPY --from=backend-builder /app/rss-reader ./rss-reader
COPY --from=frontend-builder /app/web/dist ./web/dist
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/rss-reader"]
