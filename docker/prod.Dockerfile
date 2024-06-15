FROM alpine:latest as build
ENV GO_VERSION=1.22.4

RUN apk update && apk upgrade && \
    apk add --no-cache curl tzdata

RUN update-ca-certificates

RUN curl -sL -o go.tar.gz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ ./

RUN go build -o /app/app -tags embed

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /app /
COPY .env.production /.env

ENTRYPOINT ["/app"]