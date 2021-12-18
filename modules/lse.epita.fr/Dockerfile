FROM mwendler/wget
RUN apk update && apk add bash
WORKDIR /output

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh
