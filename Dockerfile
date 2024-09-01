FROM --platform=$BUILDPLATFORM golang:1.23.0-alpine3.20@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4 AS builder

WORKDIR /usr/local/src/kube-pod-autocomplete

ENV CGO_ENABLED=0

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/kube-pod-autocomplete .


FROM alpine:3.20.2@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

COPY --from=builder /usr/local/bin/kube-pod-autocomplete /usr/local/bin/kube-pod-autocomplete

EXPOSE 8080

ENTRYPOINT [ "kube-pod-autocomplete" ]
