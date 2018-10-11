#!/usr/bin/env bash
SYS_GOCMD=$(which go)
GOBAK="$SYS_GOCMD.bak"

NEIGHBOR_DIR=$(pwd)
NEIGHBOR_GOCMD="$NEIGHBOR_DIR/bin/go-cover"

# will need to be run as `root`
if [ -f $GOBAK ]; then
  echo "Not backing up go command, a backup already exists ($GOBAK)."
else
  echo "Backing up $SYS_GOCMD to $GOBAK."
  cp $SYS_GOCMD $GOBAK
  echo "Replacing $SYS_GOCMD with $NEIGHBOR_GOCMD."
	cp $NEIGHBOR_GOCMD $SYS_GOCMD
fi
