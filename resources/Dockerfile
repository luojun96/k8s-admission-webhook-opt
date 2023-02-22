FROM golang:1.20.1-alpine as builder
WORKDIR /go/src/github.com/luojun/k8s-admission-webhook-opt
COPY ./ /go/src/github.cim/luojun/k8s-admission-webhook-opt
RUN --mount=type=cache,target=/var/cache/apk \
    apk --no-cache add make git gcc libc-dev
COPY go.mod go.mod
COPY go.sum go.sum
COPY . .
RUN --mount=type=cache,mode=0777,id=gomodcache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make

FROM alpine:latest
RUN apk update && apk add tini ca-certificates 
WORKDIR /opt/
COPY --from=builder /go/src/github.com/luojun/k8s-admission-webhook-opt/build/bin/* ./
COPY --from=builder /go/src/github.com/luojun/k8s-admission-webhook-opt/resources/scripts/entrypoint.sh /
ENV PATH /opt/:$PATH
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
