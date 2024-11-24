FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o music-service ./cmd/main.go

CMD ["./music-service"]