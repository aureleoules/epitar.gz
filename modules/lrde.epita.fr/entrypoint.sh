url="https://www.lrde.epita.fr"

# Do not download if folder already exists 
[ -d '/output/www.lrde.epita.fr' ] || \
for p in "" "~akim"; do
    wget \
        -nc \
        --no-check-certificate \
        -r \
        -np \
        --no-http-keep-alive \
        -e robots=off \
        --content-disposition \
        --reject="index.html*" \
        --reject="index.php*" \
        "$url/$p"
done

find -name 'index.php*' -delete
find -name 'index.html*' -delete

ls /output
cd /output/www.lrde.epita.fr && \
    find * \
    -type f \
    -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://www.lrde.epita.fr/$f" > "/output/www.lrde.epita.fr/$f.url"' X {} \;