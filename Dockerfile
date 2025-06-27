# =============================================================================
# 🏗️  STAGE 1 – BUILD STATIC BINARY (Go 1.23, CGO disabled)
# =============================================================================
FROM --platform=$TARGETPLATFORM golang:1.23 AS builder

WORKDIR /src

# Copy module files and download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source tree
COPY . .

# Build with a VERSION arg (defaults to “dev”)
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w -X 'main.appVersion=${VERSION}'" \
    -o /usr/local/bin/google-search-agent \
    .



# =============================================================================
# 📦  STAGE 2 – ALPINE RUNTIME (small + CA certs for HTTPS)
# =============================================================================
FROM alpine:3.18

# Install CA certs so http.Get() trusts Google
RUN apk add --no-cache ca-certificates

# Copy our binary into the image
COPY --from=builder /usr/local/bin/google-search-agent /usr/local/bin/google-search-agent

# Set working directory
WORKDIR /

# Expose the port the agent listens on
EXPOSE 8080

# At runtime, load GOOGLE_API_KEY and GOOGLE_CSE_ID from your .env file:
# docker run --rm --env-file .env -p 8080:8080 google-search-agent:latest
ENTRYPOINT ["/usr/local/bin/google-search-agent"]
CMD ["-port=8080", "-listen=0.0.0.0"]
