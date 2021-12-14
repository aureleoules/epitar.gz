FROM mwendler/wget
WORKDIR /output

ENV url "https://www.lrde.epita.fr/~akim"

ENTRYPOINT wget -nc --no-check-certificate -r -np --no-http-keep-alive --content-disposition -A.pdf ${url}
