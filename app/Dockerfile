# ---- Build Stage ----
FROM golang:1.24.2-bullseye AS builder

WORKDIR /src

# Copy go module files and download dependencies first to leverage Docker cache
COPY go.mod .
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notesapp .

# ---- Runtime Stage ----
# Use a minimal image like alpine
FROM alpine:latest

WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/notesapp .

# Environment variable for the port (optional, can be overridden)
ENV PORT=8080

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./notesapp"]
