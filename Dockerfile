FROM golang:1.17-alpine3.15 AS builder
ARG VERSION=docker-unknown
RUN apk add --no-cache make bash
COPY . /app
WORKDIR /app
RUN make build

FROM alpine:3.15
RUN apk add --no-cache bash
COPY --from=builder /app/approve-bot /approve-bot
ENTRYPOINT ["/approve-bot"]
