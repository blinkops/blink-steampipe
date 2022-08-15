#!/bin/sh

exec /home/steampipe/bin/generate && . /home/steampipe/.env && steampipe "$@"
