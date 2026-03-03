FROM golang:1.25

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

COPY testdata/ ./testdata/

RUN CGO_ENABLED=0 GOOS=linux go build -o /ovn-topo-api

EXPOSE 8080

CMD ["/ovn-topo-api"]
