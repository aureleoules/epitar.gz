FROM alpine
RUN apk update && apk add wget bash recode curl
WORKDIR /output

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh
