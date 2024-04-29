FROM golang:latest AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/dwimc ./cmd/dwimc/


FROM alpine:3.19.1

RUN mkdir /app

WORKDIR /app

COPY --from=builder /build/out/dwimc .

ENV DATABASE_URI="mongodb://mongo" \
    DATABASE_NAME="dwimc" \
    SECRET_API_KEY="please_change_me_api_key" \
    PORT="1337" \
    GIN_MODE=release

CMD [ "/app/dwimc" ]
