FROM golang:1.15 as builder

# install nodejs
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get update -y && apt-get install -y nodejs

# Build frontend part
WORKDIR /web
COPY web/*.json ./
COPY web/*.js ./
COPY web/public ./public
COPY web/src ./src
RUN npm install
RUN npm run build

# Build backend part
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY vendor ./vendor
COPY go.mod ./
COPY go.sum ./
RUN CGO_ENABLED=0 go build -mod=vendor -o poker ./cmd/poker/main.go

FROM alpine:3.11
COPY --from=builder /web/dist web/dist
COPY --from=builder /app/poker .
CMD ["./poker"]
