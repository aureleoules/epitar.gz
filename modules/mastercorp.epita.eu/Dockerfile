FROM mwendler/wget
WORKDIR /output

ENV url https://mastercorp.epita.eu/

ENTRYPOINT wget -nc --no-check-certificate -r -np -R "index.html/*" ${url}
