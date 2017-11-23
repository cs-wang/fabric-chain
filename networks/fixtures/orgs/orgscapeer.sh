#!/bin/bash

#set -x

shift $((OPTIND-1))
UP_DOWN="$@"

: ${UP_DOWN:="restart"}

. ../configall.sh

function main() {
 for org in ${ORGS[@]}; do
  cd $org
  cd ca${org##org}
   ./network_setup.sh ${UP_DOWN}
  cd ../
  for peer in ${PEERS[@]}; do
   cd $peer
    ./network_setup.sh ${UP_DOWN}
   cd ../
  done
  cd ../
 done
}

main


