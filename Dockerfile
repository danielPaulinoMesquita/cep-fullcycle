FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cep-app-fullcycle .

FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/my-go-app .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./cep-app-fullcycle"]