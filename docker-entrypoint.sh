#!/usr/bin/env bash
touch /home/steampipe/.env

/home/steampipe/blink/bin/generate

if [ $? -ne 0 ]; then
  exit $?
fi

. /home/steampipe/.env

set -Eeo pipefail

chown steampipe:0 /home/steampipe/.steampipe/db/14.2.0/data/

# if first arg is anything other than `steampipe`, assume we want to run steampipe
# this is for when other commands are passed to the container
if [ "${1:0}" != 'steampipe' ]; then
    set -- steampipe "$@"
fi

exec "$@"
