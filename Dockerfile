# syntax=docker/dockerfile:1

FROM golang:1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch AS runtime
COPY --from=builder /app/postgresql-split /postgresql-split
ENTRYPOINT ["/postgresql-split"]
