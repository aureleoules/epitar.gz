FROM mwendler/wget
run apk update && apk add bash
WORKDIR /output

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh