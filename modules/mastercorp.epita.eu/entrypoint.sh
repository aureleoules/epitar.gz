url="https://mastercorp.epita.eu"

# Do not download if folder already exists 
[ -d '/output/mastercorp.epita.eu' ] || \
wget -nc --no-check-certificate -r -np -R "index.html/*" ${url}

cd /output/mastercorp.epita.eu && \
    find * \
    -type f \
    -name "*.pdf" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://mastercorp.epita.eu/$f" > "/output/mastercorp.epita.eu/$f.url"' X {} \;