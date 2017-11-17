#!/bin/bash

#set -x

shift $((OPTIND-1))
UP_DOWN="$@"

: ${UP_DOWN:="restart"}

function main() {
 cd orderers
  bash ./orderers.sh ${UP_DOWN}
 cd ../
 cd orgs
  bash ./orgscapeer.sh ${UP_DOWN}
 cd ../
}

if [ "$UP_DOWN" = "up" ]; then
 docker network create paic
fi
main $1
if [ "$UP_DOWN" = "down" ]; then
 docker network rm paic
fi
