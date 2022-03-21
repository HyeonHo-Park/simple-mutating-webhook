# ======================
#  GO BUILD STAGE
# ======================
FROM golang:1.17-alpine3.15 as builder
WORKDIR $GOPATH/src/github.com/HyeonHo-Park/simple-mutating-webhook

ARG VERSION

COPY go.mod go.sum ./
RUN go mod verify

COPY .git .git
COPY cmd cmd
COPY internal internal

ENV GO111MODULE="on" \
  GOOS="linux" \
  CGO_ENABLED="0"

RUN apk add --no-cache \
      make \
      git && \
    rm -rf /var/cache/apk/*

RUN go install $GOPATH/src/github.com/HyeonHo-Park/simple-mutating-webhook/cmd/simple-mutating-webhook

# ======================
#  GO API STAGE
# ======================
FROM alpine:3.15
WORKDIR /simple-mutating-webhook

RUN apk add --no-cache curl && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/simple-mutating-webhook ./simple-mutating-webhook

EXPOSE 8080
CMD ["./simple-mutating-webhook"]