# Stage 1: Build
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /bit .

# Stage 2: Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /bit /usr/local/bin/bit

ENTRYPOINT ["bit"]
