FROM golang:1.19-bullseye as dev

RUN go install github.com/cosmtrek/air@v1.29.0

RUN mkdir /go/src/app

WORKDIR /go/src/app

ADD . /go/src/app

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]