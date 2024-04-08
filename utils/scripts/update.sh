#!/bin/bash

cd /app/3n-articles/ && git pull
cd /app/3n-app/ && git pull

cd /app/3n-app/ && npm run build


/app/3n-server.bin transcompile --src /app/3n-articles/markdown --dst /app/3n-articles/html
/app/3n-server.bin update -a /app/3n-articles/html