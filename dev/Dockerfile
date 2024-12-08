FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY ./ ./
RUN go mod download
RUN go build ./cmd/server
EXPOSE 8080
ENTRYPOINT ["./server"]
