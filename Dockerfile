FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go build -o kube-pod-autocomplete

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kube-pod-autocomplete .

EXPOSE 8080

CMD ["./kube-pod-autocomplete"]
