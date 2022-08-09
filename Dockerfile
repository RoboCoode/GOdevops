# syntax=docker/dockerfile:1
##build stage


FROM golang:1.19.0-alpine3.16 as base 

ENV APP_HOME /go/src/app

WORKDIR "$APP_HOME"

COPY . .

RUN go mod download


RUN go build -o app





##deploy stage



FROM alpine:latest

ENV APP_HOME /go/src/app 

RUN mkdir -p "$APP_HOME" 

WORKDIR "$APP_HOME"


COPY public public/

COPY --from=base "$APP_HOME"/app $APP_HOME

EXPOSE 3000

CMD ["./app"]
