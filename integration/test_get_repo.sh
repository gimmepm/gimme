#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath "$0")")
SCRIPT_NAME=$(basename "$0")

echo "$(date) - Starting integration test $SCRIPT_NAME"

GIMME_BIN="$SCRIPT_PATH/../bin/gimme"
COMMAND="get repos"

# the get repos command should return a zero exit status
if eval "$GIMME_BIN $COMMAND"; then
    echo "$(date) - Success at '$GIMME_BIN $COMMAND'"
else
    echo "$(date) - Failed at '$GIMME_BIN $COMMAND'"
    exit 1
fi

# get repos should return more than zero repos
REPO_COUNT=$(eval "$GIMME_BIN $COMMAND" | wc -l)
EXPECTED_REPO_COUNT=2
if [[ $REPO_COUNT -lt $EXPECTED_REPO_COUNT ]]; then
    echo "$(date) - '$GIMME_BIN $COMMAND' returned less than $EXPECTED_REPO_COUNT line(s) of output"
    exit 1
fi

echo "$(date) - All tests passed in $SCRIPT_NAME"
exit 0
