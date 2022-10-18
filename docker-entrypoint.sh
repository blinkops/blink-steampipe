#!/usr/bin/env bash
set -Eeo pipefail

chown steampipe:0 /home/steampipe/.steampipe/db/14.2.0/data/

# if first arg is anything other than `steampipe`, assume we want to run steampipe
# this is for when other commands are passed to the container
if [ "${1:0}" != 'steampipe' ]; then
    set -- steampipe "$@"
fi

if [[ $(/home/steampipe/bin/generate) -ne 0 ]]; then
  exit 1
fi

exec "$@"
