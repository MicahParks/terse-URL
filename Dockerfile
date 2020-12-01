FROM golang:1.15 AS builder

# Create a working directory and copy the code into it.
WORKDIR /app
COPY . .

# Build the code.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o terseURL cmd/terse-url-server/main.go


# The actual container.
FROM alpine

# Set some defaults for the host to bind to and the port to make it easier for people.
ENV HOST 0.0.0.0
ENV PORT 30000

# Copy the executable from the builder container.
COPY --from=builder /app/terseURL /terseURL
CMD ["/terseURL"]
