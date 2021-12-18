#!/bin/sh

mkdir -p /root/.ssh
ssh-keyscan gitlab.cri.epita.fr >> /root/.ssh/known_hosts
echo "$SSH_KEY" > /root/.ssh/id_ed25519
chmod 600 /root/.ssh/id_ed25519 
ssh-keygen -y -f /root/.ssh/id_ed25519 > /root/.ssh/id_ed25519.pub && echo "SSH key generated" || echo "SSH key generation failed"

cd epita-2025-annales && git pull origin master || git clone "git@gitlab.cri.epita.fr:danae.danycan/epita-2025-annales.git"

cd /output/epita-2025-annales && \
    find * \
    -type f \
    -name "*.pdf" \
    -exec /bin/bash -c 'f=$(printf "%s" "$1"); echo "https://gitlab.cri.epita.fr/danae.danycan/epita-2025-annales/-/raw/master/$f?inline=true" > "/output/epita-2025-annales/$f.url"' X {} \;