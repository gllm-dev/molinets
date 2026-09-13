FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates

ARG VERSION=dev
ARG TARGETOS=linux
ARG TARGETARCH=amd64

ENV GOTOOLCHAIN=local

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-w -s -X github.com/gllm-dev/molinets/internal/adapters/handlers/mcp.Version=${VERSION}" -o app cmd/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/app /usr/bin/

EXPOSE 8080

ENTRYPOINT ["app"]
CMD ["http"]
