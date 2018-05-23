#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath "$0")")
SCRIPT_NAME=$(basename "$0")

echo "$(date) - Starting integration test $SCRIPT_NAME"

GIMME_BIN="$SCRIPT_PATH/../bin/gimme"
COMMAND="get updates"

# The get updates command should return a zero exit status
if eval "$GIMME_BIN $COMMAND --since=-1day"; then
    echo "$(date) - Success at '$GIMME_BIN $COMMAND'"
else
    echo "$(date) - Failed at '$GIMME_BIN $COMMAND'"
    exit 1
fi

GITHUB_REPO="kubernetes/kubernetes"

RELEASE_DATE_RAW=$(curl -sL "https://api.github.com/repos/$GITHUB_REPO/releases/latest" |
    grep published_at |
    grep -oP ': "\K[^T]*')
RELEASE_MONTH=$(echo "$RELEASE_DATE_RAW" | grep -Po '\K[^-]*' | sed -n 2p)
MONTH_BEFORE=$(((RELEASE_MONTH - 1)))
MONTH_AFTER=$(((RELEASE_MONTH + 1)))
CURRENT_MONTH=$(date -d "$D" '+%m')

SINCE_WITH=$((((CURRENT_MONTH - MONTH_BEFORE) * -1)))
SINCE_WITHOUT=$((((CURRENT_MONTH - MONTH_AFTER) * -1)))
SINCE_CMD_WITH="--since='$SINCE_WITH months'"
SINCE_CMD_WITHOUT="--since='$SINCE_WITHOUT months'"

echo "$(date) - Running '$GIMME_BIN $COMMAND' with since to include '$GITHUB_REPO'"
if ! eval "$GIMME_BIN $COMMAND $SINCE_CMD_WITH | grep '$GITHUB_REPO'"; then
    echo "$(date) - Expected '$GITHUB_REPO' to be in the output, but it was not found"
    exit 1
fi

echo "$(date) - Running '$GIMME_BIN $COMMAND' with since to exclude '$GITHUB_REPO'"
eval "$GIMME_BIN $COMMAND $SINCE_CMD_WITHOUT | grep '$GITHUB_REPO'"
if [[ $? -ne 1 ]]; then
    echo "$(date) - Expected '$GITHUB_REPO' to be in the output, but it was not found"
    exit 1
fi

echo "$(date) - All tests passed in $SCRIPT_NAME"
exit 0
