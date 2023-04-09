FROM golang:1.20.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /serveralive

ENTRYPOINT ["/serveralive"]