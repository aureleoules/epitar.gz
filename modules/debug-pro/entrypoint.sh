[ -d '/output/www.debug-pro.com' ] || \
for url in "http://www.debug-pro.com/epita/archi/s1/en/" "http://www.debug-pro.com/epita/archi/s3/fr/" "http://www.debug-pro.com/epita/archi/s3/en/" "http://www.debug-pro.com/epita/prog/s3/"; do
wget \
    -nc \
    --no-check-certificate \
    -r \
    -erobots=off \
    -np \
    --no-http-keep-alive \
    --content-disposition -A.pdf \
    "$url"
done

ls /output
cd /output/www.debug-pro.com && \
    find * \
    -type f \
    \( -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \) \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "http://www.debug-pro.com/$f" > "/output/www.debug-pro.com/$f.url"' X {} \;