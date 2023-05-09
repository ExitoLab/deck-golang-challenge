FROM golang:1.19-alpine3.16 as builder
RUN apk --no-cache add ca-certificates git

# RUN mkdir /build
WORKDIR /opt/api

ADD . /opt/api
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

RUN go build -o main .

EXPOSE 8000

CMD ["/opt/api/main"]