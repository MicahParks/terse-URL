FROM golang:1.15 AS builder

# Get the Golang dependencies for better caching.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy in the code.
COPY . .

# Build the code.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o linkPreview cmd/social_media_link_previewer/main.go


# The actual image being produced.
FROM scratch

# Copy the executable from the builder container.
COPY --from=builder /app/linkPreview /linkPreview
COPY --from=builder /app/redirect.gohtml /
CMD ["/linkPreview"]
