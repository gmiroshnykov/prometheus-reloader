FROM golang:1.14-alpine AS builder

ARG VERSION="<development build - docker>"
ARG GIT_HASH

WORKDIR /pkg

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
RUN go build \
  -ldflags="-X 'main.Version=${VERSION}' -X 'main.GitHash=${GIT_HASH}'" \
  -o /bin/prometheus-reloader \
  cmd/main.go



FROM alpine

RUN addgroup -g 1000 -S docker && \
    adduser -u 1000 -S docker -G docker
USER 1000:1000

COPY --from=builder \
  /bin/prometheus-reloader \
  /bin/prometheus-reloader

ENTRYPOINT [ "/bin/prometheus-reloader" ]
CMD [ "-v", "1" ]
