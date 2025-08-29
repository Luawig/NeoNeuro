FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server /usr/local/bin/server
ENV HTTP_ADDR=:8080 APP_ENV=prod
EXPOSE 8080
CMD ["server"]
