#!/bin/sh

cd Past-Exams && git pull origin master || git clone "https://github.com/Epidocs/Past-Exams.git"

cd /output/Past-Exams && \
    find * \
    -type f \
    -name "*.pdf" \
    -o -name "*.doc" \
    -o -name "*.docx" \
    -o -name "*.ppt" \
    -o -name "*.pptx" \
    -o -name "*.odt" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://github.com/Epidocs/Past-Exams/blob/master/$f" > "/output/Past-Exams/$f.url"' X {} \;