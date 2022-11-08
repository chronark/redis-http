FROM golang:1.19 as builder




WORKDIR /http-proxy
COPY . .
RUN go mod download
RUN go build -o http-proxy .


FROM redis:latest

COPY --from=builder /http-proxy/http-proxy /http-proxy


ENV PORT 6379


CMD ["sh", "-c", "redis-server -v & redis-server & /http-proxy & wait"]
