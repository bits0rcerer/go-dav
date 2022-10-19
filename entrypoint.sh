#!/bin/sh

chown -R 9999:9999 /data
chmod -R 700 /data

exec su godav -c "$@"