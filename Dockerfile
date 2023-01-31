FROM golang:alpine as Builder

WORKDIR /usr/local/go/src/api

RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:3

MAINTAINER @AndreySHSH

COPY --from=Builder /usr/local/go/src/api/main /

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" >  /etc/timezone

EXPOSE 8080

ENTRYPOINT ["./main"]