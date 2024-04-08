#!/bin/bash

cd /app/3n-articles && git pull
/app/3n-server.bin -c /app/3n-app/build &
cron -f & 
wait -n
exit $?