#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath "$0")")
SCRIPT_NAME=$(basename "$0")

echo "$(date) - Starting integration test $SCRIPT_NAME"

GIMME_BIN="$SCRIPT_PATH/../bin/gimme"
COMMAND="version"

# the get repos command should return a zero exit status
if eval "$GIMME_BIN $COMMAND"; then
    echo "$(date) - Success at '$GIMME_BIN $COMMAND'"
else
    echo "$(date) - Failed at '$GIMME_BIN $COMMAND'"
    exit 1
fi

echo "$(date) - All tests passed in $SCRIPT_NAME"
exit 0
