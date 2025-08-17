#!/usr/bin/env bash

# Set WWWUSER to current user ID if not already set
# This ensures file permissions work correctly in Docker containers

if [ -z "$WWWUSER" ]; then
    export WWWUSER=${UID:-$(id -u)}
    echo "Set WWWUSER=$WWWUSER (current user ID)"
else
    echo "WWWUSER already set to $WWWUSER"
fi