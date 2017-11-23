#!/bin/bash

#set -x

. ../configall.sh

function generatepeeryaml()
{
 generateyamlpeer=$1

 cd $generateyamlpeer
 replacevar $generateyamlpeer
 chmod +x network_setup.sh
 rm docker-compose-template.yaml
 cd ../
}

function generatecayaml()
{
 generateyamlca=$1

 cd $generateyamlca
 mv docker-compose-template.yaml docker-compose.yaml
 replacecavar $generateyamlca
 chmod +x network_setup.sh
 cd ../
}

function replacevar()
{
 replacevarpeer=$1
 ARCH=`uname -s | grep Darwin`
 if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
 else
  OPTS="-i"
 fi

 NETWORK_NAME=$network_name
 sed $OPTS  "s/NETWORK_NAME/${NETWORK_NAME}/g" docker-compose.yaml
 PER_LOGMODE=$peer_log_mode
 sed $OPTS  "s/PER_LOGMODE/${PER_LOGMODE}/g" docker-compose.yaml

 if [ "$net_mode" = "local" ]; then
    sed $OPTS  "s/PEER_NAME_default/paic/g" docker-compose.yaml
 fi

 PEER_NAME=$replacevarpeer
 sed $OPTS  "s/PEER_NAME/${PEER_NAME}/g" docker-compose.yaml
 ORG_NAME=$org
 sed $OPTS  "s/ORG_NAME/${ORG_NAME}/g" docker-compose.yaml
 ORG_MSP_ID="Org"${org##org}"MSP"
 sed $OPTS  "s/ORG_MSP_ID/${ORG_MSP_ID}/g" docker-compose.yaml
 PEER_IP=`eval echo '$'"${peer}"_"${org}"_"ip"`
 sed $OPTS  "s/PEER_IP/${PEER_IP}/g" docker-compose.yaml
 PEER_PORT1=`eval echo '$'"${peer}"_"${org}"_"port1"`
 sed $OPTS  "s/PEER_PORT1/${PEER_PORT1}/g" docker-compose.yaml
 PEER_PORT2=`eval echo '$'"${peer}"_"${org}"_"port2"`
 sed $OPTS  "s/PEER_PORT2/${PEER_PORT2}/g" docker-compose.yaml

 NETWORKS_D=""
 NETWORKS_C=""
 DATA_DIR=""
 if [ "$net_mode" = "local" ]; then
  NETWORKS_D="networks:\\
  paic:\\
    external:\\
      name: paic"
  NETWORKS_C="networks:\\
      - paic"
  DATA_DIR=${data_dir}
 else
  NETWORKS_D="" 
  NETWORKS_C=""
  DATA_DIR=""
 fi
 sed $OPTS  "s/NETWORKS_D/${NETWORKS_D}/g" docker-compose.yaml
 sed $OPTS  "s/NETWORKS_C/${NETWORKS_C}/g" docker-compose.yaml
 sed $OPTS  "s/DATA_DIR/${DATA_DIR}/g" docker-compose.yaml

 if [ "$net_mode" != "local" ]; then
  genextrahosts $org docker-compose.yaml $OPTS
 fi

}

function replacecavar()
{
 replacevarca=$1
 ARCH=`uname -s | grep Darwin`
 if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
 else
  OPTS="-i"
 fi

 NETWORK_NAME=$network_name
 sed $OPTS  "s/NETWORK_NAME/${NETWORK_NAME}/g" docker-compose.yaml
 CA_NAME_ID=$replacevarca
 sed $OPTS  "s/CA_NAME_ID/${CA_NAME_ID}/g" docker-compose.yaml
 CA_ORG_NAME=ca-$org
 sed $OPTS  "s/CA_ORG_NAME/${CA_ORG_NAME}/g" docker-compose.yaml
 ORG_MSP_ID="Org"${org##org}"MSP"
 sed $OPTS  "s/ORG_MSP_ID/${ORG_MSP_ID}/g" docker-compose.yaml
 ORG_NAME=$org
 sed $OPTS  "s/ORG_NAME/${ORG_NAME}/g" docker-compose.yaml

 CURRENT_DIR=$PWD
 cd ../../../blockchainkey/crypto-config/peerOrganizations/$org.$network_name.com/ca/
 CA_PRIVATE_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 sed $OPTS  "s/CA_PRIVATE_KEY/${CA_PRIVATE_KEY}/g" docker-compose.yaml
 cd ../../../blockchainkey/crypto-config/peerOrganizations/$org.$network_name.com/ca/
 TLS_PRIVATE_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 sed $OPTS  "s/TLS_PRIVATE_KEY/${TLS_PRIVATE_KEY}/g" docker-compose.yaml

 CA_PORT=`eval echo '$'"ca${org##org}"_"${org}"_"port"`
 sed $OPTS  "s/CA_PORT/${CA_PORT}/g" docker-compose.yaml

 NETWORKS_D=""
 NETWORKS_C=""
 DATA_DIR=""
 if [ "$net_mode" = "local" ]; then
  NETWORKS_D="networks:\\
  paic:\\
    external:\\
      name: paic"
  NETWORKS_C="networks:\\
      - paic"
  DATA_DIR=${data_dir}
 else
  NETWORKS_D="" 
  NETWORKS_C="" 
 fi
 sed $OPTS  "s/NETWORKS_D/${NETWORKS_D}/g" docker-compose.yaml
 sed $OPTS  "s/NETWORKS_C/${NETWORKS_C}/g" docker-compose.yaml
 sed $OPTS  "s/DATA_DIR/${DATA_DIR}/g" docker-compose.yaml

}


function main() {
 if [ "$1" = "" ]; then 
  for org in ${ORGS[@]}; do
   cp -r org-template $org
   cd $org 
   mv ca-template ca${org##org}
   generatecayaml ca${org##org}
   for peer in ${PEERS[@]}; do
    echo $peer.$org.$network_name.com
    cp -r peer-template $peer
    generatepeeryaml $peer
   done
   rm -rf peer-template
   cd ..
  done
 else
  for org in ${ORGS[@]}; do
   sudo rm -rf $org
  done
 fi
}

main $1


