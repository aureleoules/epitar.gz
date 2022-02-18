FROM alpine:latest
RUN apk add git bash openssh

WORKDIR /output

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT /entrypoint.sh