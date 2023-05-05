FROM golang:1.19.2-alpine3.16 as builder

WORKDIR /go/src/github.com/ebogdanov/su27bot/
COPY . /go/src/github.com/ebogdanov/su27bot/

RUN go mod tidy -v

RUN CGO_ENABLED=0 go build -o agent main.go

FROM alpine:latest

WORKDIR /opt/su27bot
COPY --from=builder /go/src/github.com/ebogdanov/su27bot/agent ./
# COPY conf/config.yaml ./conf/config.yaml

CMD ["./agent"]
