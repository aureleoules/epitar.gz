FROM mwendler/wget
WORKDIR /output

ENV url "https://www.lse.epita.fr/"

ENTRYPOINT wget --no-check-certificate -r -np --no-http-keep-alive --content-disposition -A.pdf ${url}
