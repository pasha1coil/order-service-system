FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/notification-service ./notification_service/cmd

FROM gcr.io/distroless/base-debian12
WORKDIR /srv
COPY --from=builder /out/notification-service /srv/notification-service
ENTRYPOINT ["/srv/notification-service"]


