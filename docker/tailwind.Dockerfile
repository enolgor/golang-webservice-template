FROM --platform=$BUILDPLATFORM alpine:latest
ARG version
RUN apk update && apk upgrade && \
    apk add --no-cache curl
ARG BUILDPLATFORM
RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then \
      curl -sL -o /tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-linux-x64"; \
    elif [ "$BUILDPLATFORM" = "linux/arm64" ]; then \
    curl -sL -o /tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-linux-arm64"; \
    else \
      echo "Unsupported platform"; \
      exit 1; \
    fi
RUN chmod +x /tailwindcss
ENTRYPOINT ["/tailwindcss"]