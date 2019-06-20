#!/bin/ash
set -e
./envnginx -d=/etc/nginx/conf.d/*.conf
nginx -g "daemon off;"
