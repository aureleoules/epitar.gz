#!/bin/sh

cd epiportal.com && git pull origin master || git clone "https://github.com/Epiportal/epiportal.com.git"

cd /output/epiportal.com && \
    find * \
    -type f \
    \( -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \) \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://github.com/Epiportal/epiportal.com/blob/master/$f?inline=true" > "/output/epiportal.com/$f.url"' X {} \;