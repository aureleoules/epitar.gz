#!/bin/sh

curl \
     -c cookies \
     -v \
     -L \
     -H "Referer: https://epinotes.fr/public/pages/login_epid.html" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "mail=${EMAIL}&epid=${EPID}"  \
     -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
    "https://epinotes.fr/connection/connect_epid.php"

echo 'Downloading...'

wget \
    -nc \
    --no-check-certificate \
     --load-cookies=cookies \
    -r \
    -np \
    --no-http-keep-alive \
    -e robots=off \
    --content-disposition \
    --verbose \
    "https://epinotes.fr/documents/data/ftp.php?dir=/"

cd epinotes.fr && \
    find * \
    -type f \
    \( -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \) \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://epinotes.fr/$f" > "/output/epinotes.fr/$f.url"' X {} \;
