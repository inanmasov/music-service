FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o music-api ./main.go

CMD ["./music-api"]