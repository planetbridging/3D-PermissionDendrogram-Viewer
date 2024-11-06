# Use an official Golang image as a build environment
FROM golang:1.20 as builder

WORKDIR /app
COPY . .

# Build the Go app
RUN go mod tidy
RUN go build -o main .

# Use a minimal base image to run the application
FROM alpine:latest
WORKDIR /root/

# Copy the Go binary from the builder image
COPY --from=builder /app/main .

EXPOSE 8432

# Run the binary
CMD ["./main"]
