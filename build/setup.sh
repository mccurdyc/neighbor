#!/usr/bin/env bash
NEIGHBOR_DIR=$1
NEIGHBOR_CFG="$NEIGHBOR_DIR/config.json"
NEIGHBOR_SAMPLE_CFG="$NEIGHBOR_DIR/sample.config.json"

if [ -f $NEIGHBOR_CFG ]; then
  echo "neighbor config already exists ($NEIGHBOR_CFG), doing nothing."
else
  echo "Creating a neighbor config file ($NEIGHBOR_CFG) from the sample config."
	cp $NEIGHBOR_SAMPLE_CFG $NEIGHBOR_CFG
fi
