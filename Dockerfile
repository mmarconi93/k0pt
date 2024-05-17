# Stage 1: Build
FROM golang:latest AS builder
WORKDIR /src/app

# Copy the source code
COPY . .

# Build the application
RUN go build -o k0pt cmd/main.go 

# Stage 2: Runtime
FROM alpine:latest AS runner
WORKDIR /usr/local/bin

# Copy the binary from the builder stage
COPY --from=builder /src/app/k0pt .

# Run the binary
ENTRYPOINT ["./k0pt"]
CMD ["version"]