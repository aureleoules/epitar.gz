url="https://mastercorp.epita.eu"

# Do not download if folder already exists 
[ -d '/output/mastercorp.epita.eu' ] || \
wget -nc --no-check-certificate -r -np -R "index.html/*" ${url}

cd /output/mastercorp.epita.eu && \
    find * \
    -type f \
    -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://mastercorp.epita.eu/$f" > "/output/mastercorp.epita.eu/$f.url"' X {} \;