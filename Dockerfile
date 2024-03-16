FROM golang:latest AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/dwimc ./cmd/dwimc/


FROM alpine:3.15.0

RUN mkdir /app && mkdir /data
VOLUME [ "/data" ]

WORKDIR /app

COPY --from=builder /build/out/dwimc .

ENV DATABASE_URI="mongodb://mongo" \
    DATABASE_NAME="dwimc" \
    PORT="1337" \
    SECRET_API_KEY="" \
    GIN_MODE=release

CMD [ "/app/dwimc" ]
