FROM golang:latest AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make build

FROM alpine:3.21.2

RUN mkdir /app

WORKDIR /app

COPY --from=builder /build/bin/dwimc .

CMD [ "/app/dwimc" ]
