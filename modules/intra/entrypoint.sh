#!/bin/sh

ls /output/* && exit 0
Xvfb :99 -screen 0 1920x1080x16 & node index.js