FROM golang:1.17-alpine3.15 AS builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o ./approver-bot

FROM alpine:3.15
COPY --from=builder /app/approver-bot /approver-bot
ENTRYPOINT ["/approver-bot"]