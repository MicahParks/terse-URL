FROM golang:1.15 AS builder

# Get the Golang dependencies for better caching.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the code in.
COPY . .

# Build the code.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o terseURL cmd/terse-url-server/main.go


# The actual image being produced.
FROM scratch

# Set some defaults for the host to bind to and the port to make it easier for people.
ENV HOST 0.0.0.0
ENV PORT 30000

# Copy the executable from the builder container.
COPY --from=builder /app/terseURL /terseURL
CMD ["/terseURL"]