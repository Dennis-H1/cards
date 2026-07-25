# ---- frontend build ----
FROM node:22-alpine AS frontend-build
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- go build ----
FROM golang:1.25-alpine AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/karteikarten ./cmd/karteikarten

# ---- backend image ----
FROM alpine:3.20 AS backend
RUN apk add --no-cache ca-certificates
COPY --from=go-build /out/karteikarten /usr/local/bin/karteikarten
EXPOSE 8080
ENTRYPOINT ["karteikarten"]

# ---- web image: serves the built frontend, proxies /api, /mcp, /healthz to backend ----
FROM nginx:1.27-alpine AS web
COPY --from=frontend-build /src/frontend/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
