FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o cep-app-fullcycle .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/cep-app-fullcycle /app/

USER 65536:65536

# Optionally, expose the port
EXPOSE 8080

# Start the application
CMD ["/app/cep-app-fullcycle"]