#!/bin/bash

shift $((OPTIND-1))
UP_DOWN="$@"

: ${UP_DOWN:="restart"}

. ../configall.sh

function genallinone() {
 mkdir zkfks
  for zookeeper in ${ZOOKEEPERS[@]}; do
  cp $zookeeper/docker-compose.yaml zkfks/docker-compose-$zookeeper.yaml
 done

 for kafka in ${KAFKAS[@]}; do
  cp $kafka/docker-compose.yaml zkfks/docker-compose-$kafka.yaml
 done

}


function main() {
 for zookeeper in ${ZOOKEEPERS[@]}; do
  cd $zookeeper
  ./network_setup.sh ${UP_DOWN}
  cd ../
 done
 
 for kafka in ${KAFKAS[@]}; do
  cd $kafka
  ./network_setup.sh ${UP_DOWN}
  cd ../
 done
}

if [ "$1" != "gen" ]; then
 main
else
 genallinone
fi
