#!/bin/sh

# Support the following running mode
# 1) run webhook-server without any arguments
# 2) run the command in webhook-server container

if [ "$#" -eq 0 ]; then
   exec /opt/webhook-server
else
    exec "$@"
fi
