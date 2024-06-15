FROM alpine:latest
ENV GO_VERSION=1.22.4

RUN apk update && apk upgrade && \
    apk add --no-cache curl tzdata

RUN update-ca-certificates

RUN curl -sL -o go.tar.gz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download

ENTRYPOINT ["air", "-c", ".air.toml"]
