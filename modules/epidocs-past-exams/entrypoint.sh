#!/bin/sh

cd Past-Exams && git pull origin master || git clone "https://github.com/Epidocs/Past-Exams.git"

cd /output/Past-Exams && \
    find * \
    -type f \
    -name "*.pdf" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://github.com/Epidocs/Past-Exams/blob/master/$f" > "/output/Past-Exams/$f.url"' X {} \;