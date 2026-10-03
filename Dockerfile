FROM node:22-alpine AS frontend
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci --legacy-peer-deps
COPY web/babel.config.js web/vue.config.js ./
COPY web/public ./public
COPY web/src ./src
# Vue CLI 4 uses Webpack 4, which requires the legacy OpenSSL provider.
RUN NODE_OPTIONS=--openssl-legacy-provider npm run build

FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags="-s -w" -o /poker ./cmd/poker

FROM alpine:3.23
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 poker
WORKDIR /app
COPY --from=frontend /web/dist ./web/dist
COPY --from=backend /poker ./poker
ENV GIN_MODE=release PORT=8080
USER poker
EXPOSE 8080
CMD ["./poker"]
