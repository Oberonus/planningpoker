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

# Difference between production and dev image is that dev version does not use builder
# container which allowes to use layers caching and run local CI faster
WORKDIR /poker
RUN mkdir -p web/dist
RUN cp -R /web/dist ./web/dist
RUN cp /app/poker .
CMD ["./poker"]
