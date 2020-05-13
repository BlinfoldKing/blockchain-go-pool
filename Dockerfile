FROM alpine:3.7
WORKDIR /server
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
COPY ./bin/app .
COPY ./schema ./schema
EXPOSE 9000
CMD ["/server/app"]