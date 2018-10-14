#!/usr/bin/env bash
SYS_GOCMD=$(which go)
GOBAK="$SYS_GOCMD.bak"

NEIGHBOR_DIR=$(pwd)
NEIGHBOR_GOCMD="$NEIGHBOR_DIR/bin/go-cover"
NEIGHBOR_CFG="$NEIGHBOR_DIR/config.json"
NEIGHBOR_SAMPLE_CFG="$NEIGHBOR_DIR/sample.config.json"

if [ -f $NEIGHBOR_CFG ]; then
  echo "neighbor config already exists ($NEIGHBOR_CFG), doing nothing."
else
  echo "Creating a neighbor config file ($NEIGHBOR_CFG) from the sample config."
	cp $NEIGHBOR_SAMPLE_CFG $NEIGHBOR_CFG
fi

# will need to be run as `root`
if [ -f $GOBAK ]; then
  echo "Not backing up go command, a backup already exists ($GOBAK)."
else
  echo "Backing up $SYS_GOCMD to $GOBAK."
  cp $SYS_GOCMD $GOBAK
  echo "Replacing $SYS_GOCMD with $NEIGHBOR_GOCMD."
	cp $NEIGHBOR_GOCMD $SYS_GOCMD
fi
