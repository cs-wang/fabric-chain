#!/bin/bash

#set -x

shift $((OPTIND-1))
UP_DOWN="$@"

: ${UP_DOWN:="restart"}

. ../configall.sh

function main() {
 for orderer in ${ORDERERS[@]}; do
  cd $orderer
  ./network_setup.sh -s ${UP_DOWN}
  cd ../
 done
}

main $1
