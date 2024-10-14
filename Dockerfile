FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/mg_vault ./mg_vault
CMD ("/app/mg_vault") 