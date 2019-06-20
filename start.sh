#!/bin/ash
set -e
./nginx-config -d=/etc/nginx/conf.d/*.conf
nginx -g "daemon off;"
