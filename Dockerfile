FROM golang:1.24.1 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build


FROM alpine:3.21.3

RUN apk add --no-cache ca-certificates && \
    mkdir /app && \
    adduser -D -g '' dwimcuser && \
    chown -R dwimcuser /app

WORKDIR /app

COPY --from=builder /build/bin/dwimc .
RUN chmod +x /app/dwimc && \
    chown dwimcuser /app/dwimc

USER dwimcuser

ENV DATABASE_URI="mongodb://mongo" \
    DATABASE_NAME="dwimc" \
    PORT=1337 \
    GIN_MODE="release" \
    LOG_OUTPUT_TYPE="console" \
    LOG_LEVEL="info" \
    DEBUG_MODE=false \
    LOCATION_HISTORY_LIMIT=10

CMD [ "/app/dwimc" ]
