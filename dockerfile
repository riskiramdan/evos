# Copyright 2019 Core Services Team.

FROM golang:1.15-alpine as builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install ./cmd/evos

FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin /bin
COPY --from=builder /app/.env /bin

USER nobody:nobody
ENTRYPOINT ["/bin/evos"]