#!/bin/sh

alias urldecode='sed "s@+@ @g;s@%@\\\\x@g" | xargs -0 printf "%b"'

# Retrieve CRI's login page to get CSRF tokeen
curl \
    -c csrf \
    -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
    "https://cri.epita.fr/auth/login/" &> resp

# Retrievee CSRF token
TOKEN=$(sed -n 's:.*"csrfmiddlewaretoken" value="\(.*\)".*:\1:p' resp)
echo $TOKEN

# Login with credentials and token
curl -b csrf \
     -c csrf \
     -v \
     -L \
     -H "Referer: https://cri.epita.fr/auth/login/" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "csrfmiddlewaretoken=${TOKEN}&username=${LOGIN}&password=${PASSWORD}"  \
     -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
    "https://cri.epita.fr/auth/login/" &> resp

CRISESSIONID=$(sed -n 's:.*sessionid="\(.*\)".*:\1:p' resp | head -1)
echo "$CRISESSIONID"

# Retrieve URL with session key in HTML
curl -c csrf \
     -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
     -v \
    "https://moodle.cri.epita.fr/login/index.php" &> index.php

URL=$(sed -n 's:.*href="\(.*\)" title="CRI".*:\1:p' index.php | recode html..ascii | urldecode)
echo $URL

# Redirect flow towards CRI
curl -L -H "Cookie: sessionid=$CRISESSIONID" \
    -H "Referer: https://moodle.cri.epita.fr/login/index.php" \
    -c csrf -b csrf \
    -v \
     -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
    "$URL" &> resp

# Retrieve URL to go to
URL=$(sed -n 's:.*location\: \(.*\).*:\1:p' resp | head -1)
URL=${URL::-1}

echo "$URL"

curl -L -H "cookie: sessionid=$CRISESSIONID" \
    -H "Referer: https://moodle.cri.epita.fr/login/index.php" \
    -c csrf -b csrf \
    -v \
     -A "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
    "$URL" &> resp

# Archive with cookies
wget \
    -nc \
    --no-check-certificate \
     --load-cookies=csrf \
    -r \
    -np \
    --no-http-keep-alive \
    -e robots=off \
    --content-disposition \
    --verbose \
    --reject-regex logout \
    "https://moodle.cri.epita.fr"


cd moodle.cri.epita.fr && \
    find * \
    -type f \
    -name "*.pdf" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://moodle.cri.epita.fr/$f" > "/output/moodle.cri.epita.fr/$f.url"' X {} \;