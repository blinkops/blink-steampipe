#!/bin/sh

touch /home/steampipe/.env

exec /home/steampipe/bin/generate && . /home/steampipe/.env && steampipe "$@"
