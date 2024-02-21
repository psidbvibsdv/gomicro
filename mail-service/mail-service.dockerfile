FROM golang:1.21-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api


#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailerApp /app
COPY --from=builder /app/templates ../templates

CMD ["/app/mailerApp"]