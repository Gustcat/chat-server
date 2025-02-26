FROM golang:1.23.2-alpine AS builder

COPY . /github.com/Gustcat/chat-server/source/
WORKDIR /github.com/Gustcat/chat-server/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Gustcat/chat-server/source/bin/crud_server .
COPY --from=builder /github.com/Gustcat/chat-server/source/.env .

CMD ["./crud_server"]