#!/bin/bash

#set -x

. ../configall.sh

function generateyaml()
{
 generateyamlorderer=$1
 cd $generateyamlorderer
 replacevar $generateyamlorderer
 chmod +x network_setup.sh
 rm docker-compose-template.yaml
 cd ../
}

function replacevar()
{
 replacevarorderer=$1
 ARCH=`uname -s | grep Darwin`
 if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
 else
  OPTS="-i"
 fi

 NETWORK_NAME=$network_name
 sed $OPTS  "s/NETWORK_NAME/${NETWORK_NAME}/g" docker-compose.yaml
 ORDERER_NAME=$replacevarorderer
 sed $OPTS  "s/ORDERER_NAME/${ORDERER_NAME}/g" docker-compose.yaml
 ORDERER_PORT=`eval echo '$'"${orderer}"_"port"`
 sed $OPTS  "s/ORDERER_PORT/${ORDERER_PORT}/g" docker-compose.yaml
 ORDERER_IP=`eval echo '$'"${orderer}"_"ip"`
 sed $OPTS  "s/ORDERER_IP/${ORDERER_IP}/g" docker-compose.yaml

 genextrahosts $orderer docker-compose.yaml $OPTS
}

function main() {
 if [ "$1" = "" ]; then 
  for orderer in ${ORDERERS[@]}; do
   echo "$orderer".$network_name.com
   cp -r orderer-template $orderer
   generateyaml $orderer
  done
 else
  for orderer in ${ORDERERS[@]}; do
   rm -rf $orderer
  done
  sudo rm -rf ./data
 fi
}

main $1
