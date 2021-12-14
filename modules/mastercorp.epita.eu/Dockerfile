FROM mwendler/wget
WORKDIR /output

ENV url https://mastercorp.epita.eu/

ENTRYPOINT wget --no-check-certificate -r -np -R "index.html/*" ${url}
