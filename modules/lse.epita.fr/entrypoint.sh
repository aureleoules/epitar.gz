url="https://www.lse.epita.fr"

# Do not download if folder already exists 
[ -d '/output/www.lse.epita.fr' ] || \
wget -nc --no-check-certificate -r -np --no-http-keep-alive --content-disposition -A.pdf ${url}

cd /output/www.lse.epita.fr && \
    find * \
    -type f \
    -name "*.pdf" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://www.lse.epita.fr/$f" > "/output/www.lse.epita.fr/$f.url"' X {} \;