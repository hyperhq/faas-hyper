FROM golang:1.7.5

RUN mkdir -p /go/src/github.com/hyperhq/faas-hyper/

WORKDIR /go/src/github.com/hyperhq/faas-hyper

COPY vendor vendor
COPY handlers	handlers
COPY hyper hyper
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faas-hyper .

EXPOSE 8080

CMD ["./faas-hyper"]
