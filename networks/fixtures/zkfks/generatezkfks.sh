#!/bin/bash

#set -x

. ../configall.sh

function generatezookeeperyaml()
{
 generateyamlzookeeper=$1

 cd $generateyamlzookeeper
 replacezookeepervar $generateyamlzookeeper
 chmod +x network_setup.sh
 rm docker-compose-template.yaml
 cd ../
}

function generatekafkayaml()
{
 generateyamlkafka=$1
 cd $generateyamlkafka
 replacekafkavar $generateyamlkafka
 chmod +x network_setup.sh
 rm docker-compose-template.yaml
 cd ../
}

function replacezookeepervar()
{

 replacezookeepervar=$1
 ARCH=`uname -s | grep Darwin`
 if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
 else
  OPTS="-i"
 fi

 ZOOKEEPER_NAME=$replacezookeepervar
 sed $OPTS  "s/ZOOKEEPER_NAME/${ZOOKEEPER_NAME}/g" docker-compose.yaml
 ZOOKEEPER_ID=$((${replacezookeepervar##zookeeper}+1))
 sed $OPTS  "s/ZOOKEEPER_ID/${ZOOKEEPER_ID}/g" docker-compose.yaml

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
  NETWORKS_C="networks_mode: host\\
ports:\\
  - 2181:2181\\
  - 2888:2888\\
  - 3888:3888"
 fi
 sed $OPTS  "s/NETWORKS_D/${NETWORKS_D}/g" docker-compose.yaml 
 sed $OPTS  "s/NETWORKS_C/${NETWORKS_C}/g" docker-compose.yaml 
 sed $OPTS  "s/DATA_DIR/${DATA_DIR}/g" docker-compose.yaml 

 if [ "$net_mode" != "local" ]; then
  genextrahosts $zookeeper docker-compose.yaml $OPTS
 fi

}

function replacekafkavar()
{
 replacekafkavar=$1
 ARCH=`uname -s | grep Darwin`
 if [ "$ARCH" == "Darwin" ]; then
  OPTS="-it"
 else
  OPTS="-i"
 fi

 KAFKA_NAME=$replacekafkavar
 sed $OPTS  "s/KAFKA_NAME/${KAFKA_NAME}/g" docker-compose.yaml
 KFK_BK_ID=$((${replacekafkavar##kafka}+1))
 sed $OPTS  "s/KFK_BK_ID/${KFK_BK_ID}/g" docker-compose.yaml
 KFK_ADV_HN=`eval echo '$'"${replacekafkavar}"_"ip"`
 sed $OPTS  "s/KFK_ADV_HN/${KFK_ADV_HN}/g" docker-compose.yaml
 KFK_ADV_PT=`eval echo '$'"${replacekafkavar}"_"port"`
 sed $OPTS  "s/KFK_ADV_PT/${KFK_ADV_PT}/g" docker-compose.yaml
 KAFKA_PORT=`eval echo '$'"${replacekafkavar}"_"port"`
 sed $OPTS  "s/KAFKA_PORT/${KAFKA_PORT}/g" docker-compose.yaml

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
  NETWORKS_C="networks_mode: host"
 fi
 sed $OPTS  "s/NETWORKS_D/${NETWORKS_D}/g" docker-compose.yaml
 sed $OPTS  "s/NETWORKS_C/${NETWORKS_C}/g" docker-compose.yaml 
 sed $OPTS  "s/DATA_DIR/${DATA_DIR}/g" docker-compose.yaml 

 if [ "$net_mode" != "local" ]; then
  genextrahosts $kafka docker-compose.yaml $OPTS
 fi

}


function main() {
 if [ "$1" = "" ]; then 
  for zookeeper in ${ZOOKEEPERS[@]}; do
   echo $zookeeper
   cp -r zookeeper-template $zookeeper
   generatezookeeperyaml $zookeeper
  done
  
  for kafka in ${KAFKAS[@]}; do
   echo $kafka
   cp -r kafka-template $kafka
   generatekafkayaml $kafka
  done
 else
  for zookeeper in ${ZOOKEEPERS[@]}; do
   sudo rm -rf $zookeeper
  done
  for kafka in ${KAFKAS[@]}; do
   sudo rm -rf $kafka
  done
 fi
}

main $1


