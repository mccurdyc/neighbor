#!/usr/bin/env bash
SYS_GOCMD=$(which go)
GOBAK="$SYS_GOCMD.bak"

NEIGHBOR_DIR=$(pwd)
NEIGHBOR_GOCMD="$NEIGHBOR_DIR/bin/go-cover"

# will need to be run as `root`
if [ -f $GOBAK ]; then
  echo "Moving backup $GOBAK back to $SYS_GOCMD."
  sudo mv $GOBAK $SYS_GOCMD
else
  echo "Backup ($GOBAK) doesn't exist; doing nothing."
fi
