#!/bin/sh


[ -d '/output/annales.hyperion.tf' ] && exit 0

if [ -z "${USERNAME}" ]; then
  echo "USERNAME is not set"
  exit 1
fi

# check if USERPASSWORD is set
if [ -z "${PASSWORD}" ]; then
  echo "PASSWORD is not set"
  exit 1
fi

wget --keep-session-cookies \
     --save-cookies=csrf.txt \
    --no-check-certificate \
     -U "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
      "https://annales.hyperion.tf/login/"

TOKEN=$(sed -n 's:.*"csrfmiddlewaretoken" value="\(.*\)".*:\1:p' index.html)
echo $TOKEN

wget -v --keep-session-cookies \
     --load-cookies=csrf.txt \
    --no-check-certificate \
     --save-cookies=cookies.txt \
     -U "Mozilla/5.0 (X11; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0" \
     -d \
     --header "Referer: https://annales.hyperion.tf/login/" \
     --post-data="csrfmiddlewaretoken=${TOKEN}&username=${USERNAME}&password=${PASSWORD}"  \
     "https://annales.hyperion.tf/login/"


cat cookies.txt
sleep 1

wget \
    -nc \
    --no-check-certificate \
    --load-cookies=cookies.txt \
    -r \
    -np \
    --no-http-keep-alive \
    -e robots=off \
    "https://annales.hyperion.tf"

cd annales.hyperion.tf
for d in $( find . -type d -print ) ; do
   if [ -f $d/index.html ] ; then
       mv $d/index.html $d/$(basename $d)
    fi
done


find * \
    -type f \
    \( -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \) \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); d=$(dirname "$f"); echo "https://annales.hyperion.tf/$d" > "/output/annales.hyperion.tf/$f.url"' X {} \;
