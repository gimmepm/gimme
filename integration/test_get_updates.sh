#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath "$0")")

GIMME_BIN="$SCRIPT_PATH/../bin/gimme"
COMMAND="get updates"

# the get updates command should return a zero exit status
if eval "$GIMME_BIN $COMMAND"; then
    echo "$(date) - Success at '$GIMME_BIN $COMMAND'"
else
    echo "$(date) - Failed at '$GIMME_BIN $COMMAND'"
    exit 1
fi

# get updates should return more than zero updates
UPDATE_OUTPUT=$(eval "$GIMME_BIN $COMMAND")
UPDATE_OUTPUT_COUNT=$(echo "$UPDATE_OUTPUT" | wc -l)
EXPECTED_UPDATE_COUNT=2
if [[ $UPDATE_OUTPUT_COUNT -lt $EXPECTED_UPDATE_COUNT ]]; then
    echo "$(date) - '$GIMME_BIN $COMMAND' returned less than $EXPECTED_UPDATE_COUNT line(s) of output"
    echo "Output:"
    echo "$UPDATE_OUTPUT"
    exit 1
fi

echo "$(date) - All tests passed"
exit 0
