FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/billing-service ./billing_service/cmd

FROM gcr.io/distroless/base-debian12
WORKDIR /srv
COPY --from=builder /out/billing-service /srv/billing-service
ENTRYPOINT ["/srv/billing-service"]


