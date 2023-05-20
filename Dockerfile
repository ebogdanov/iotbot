FROM golang:1.19.2-alpine as builder

RUN apk update && apk add sqlite && apk add --update gcc musl-dev

WORKDIR /go/src/github.com/ebogdanov/su27bot/
COPY . /go/src/github.com/ebogdanov/su27bot/

RUN go mod tidy -v
RUN CGO_ENABLED=1 go build -tags "linux" -o agent main.go

FROM alpine:latest
WORKDIR /opt/su27bot
COPY --from=builder /go/src/github.com/ebogdanov/su27bot/agent ./
# COPY conf/config.yaml ./conf/config.yaml

CMD ["./agent"]
