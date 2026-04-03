FROM golang:tip-alpine3.23 AS builder

WORKDIR /app
COPY main.go .
RUN go build -o server main.go

FROM scratch
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
