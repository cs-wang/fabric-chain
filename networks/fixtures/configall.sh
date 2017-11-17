#set -x
#########################################
# prehistoric 
#########################################
extra_host="extra_hosts:"
blank="      "
host_end="extra_host_ip"

orderer_template_path=orderers/orderer-template
peer_template_path=orgs/org-template/peer-template
zookeeper_template_path=zkfks/zookeeper-template
kafka_template_path=zkfks/kafka-template

src_template=docker-compose-template.yaml
dest_template=docker-compose.yaml

network_name=example

data_path=/data

#net_mode=local

LAN_ORG_UNITS=(a_unit)
a_unit=(zookeeper0 zookeeper1 zookeeper2 kafka0 kafka1 kafka2 orderer0 orderer1 orderer2 org1 org2 org3 org4)

#########################################
# orderers
#########################################
ORDERERS=(orderer0 orderer1 orderer2)
orderer0_ip=
orderer0_local_ip=10.1.2.62
orderer0_port=7050

orderer1_ip=
orderer1_local_ip=10.1.2.60
orderer1_port=7050

orderer2_ip=
orderer2_local_ip=10.1.2.61
orderer2_port=7050

#########################################
# orgs
#########################################
ORGS=(org1 org2 org3 org4)
PEERS=(peer0 peer1)

ca1_org1_ip=
ca1_org1_port=7054
ca2_org2_ip=
ca2_org2_port=7054
ca3_org3_ip=
ca3_org3_port=7054
ca4_org4_ip=
ca4_org4_port=7054


peer0_org1_ip=
peer0_org1_local_ip=10.1.17.149
peer0_org1_port1=7051
peer0_org1_port2=7053

peer1_org1_ip=
peer1_org1_local_ip=10.1.36.208
peer1_org1_port1=8051
peer1_org1_port2=8053

peer0_org2_ip=
peer0_org2_local_ip=10.1.36.208
peer0_org2_port1=7051
peer0_org2_port2=7053

peer1_org2_ip=
peer1_org2_local_ip=10.1.17.149
peer1_org2_port1=8051
peer1_org2_port2=8053

peer0_org3_ip=
peer0_org3_local_ip=10.4.4.61
peer0_org3_port1=7051
peer0_org3_port2=7053

peer1_org3_ip=
peer1_org3_local_ip=10.4.4.62
peer1_org3_port1=7051
peer1_org3_port2=7053

peer0_org4_ip=
peer0_org4_local_ip=10.4.4.64
peer0_org4_port1=7051
peer0_org4_port2=7053

peer1_org4_ip=
peer1_org4_local_ip=10.4.4.65
peer1_org4_port1=7051
peer1_org4_port2=7053
#########################################
# zkfks
#########################################
ZOOKEEPERS=(zookeeper0 zookeeper1 zookeeper2)
KAFKAS=(kafka0 kafka1 kafka2)

zookeeper0_ip=
zookeeper0_local_ip=
zookeeper1_ip=
zookeeper1_local_ip=
zookeeper2_ip=
zookeeper2_local_ip=

kafka0_ip=
kafka0_local_ip=10.1.2.62
kafka0_port=9092

kafka1_ip=
kafka1_local_ip=10.1.2.60
kafka1_port=9092

kafka2_ip=
kafka2_local_ip=10.1.2.61
kafka2_port=9092

#########################################
################extra_host###############
#########################################
function getorgunit() {
 utype=$1
 for unit in ${LAN_ORG_UNITS[@]}; do
  for orga in $(eval echo '${'"${unit}""[@]}"); do
   if [ "$utype" = "$orga" ]; then
     echo $unit
     return 0
   fi
  done
 done
 return 1
}

function genextrahosts() {
 type=$1
 file=$2
 OPTS=$3
 typeorg=$(getorgunit $type)
 ipforreplace=""

 for orderer in ${ORDERERS[@]}; do
    theorg=$(getorgunit $orderer)
    orderer_ip="${orderer}_${host_end}"
    if [ "$typeorg" = "$theorg" ]; then
      ipforreplace=`eval echo '$'"${orderer}"_"local_ip"`
    else
      ipforreplace=`eval echo '$'"${orderer}"_"ip"`
    fi
    sed $OPTS  "s/${orderer_ip}/${ipforreplace}/g" $file
 done

 for orgb in ${ORGS[@]}; do
   theorg=$(getorgunit $orgb)
   for peer in ${PEERS[@]}; do
     peer_ip="${peer}_${orgb}_${host_end}"
     if [ "$typeorg" = "$theorg" ]; then
      ipforreplace=`eval echo '$'"${peer}"_"${orgb}"_"local_ip"`
     else
      ipforreplace=`eval echo '$'"${peer}"_"${orgb}"_"ip"`
     fi
     sed $OPTS  "s/${peer_ip}/${ipforreplace}/g" $file
   done
 done

 for kafka in ${KAFKAS[@]}; do
   theorg=$(getorgunit $kafka)
   kafka_ip="${kafka}_${host_end}"
   if [ "$typeorg" = "$theorg" ]; then
     ipforreplace=`eval echo '$'"${kafka}"_"local_ip"`
   else
     ipforreplace=`eval echo '$'"${kafka}"_"ip"`
   fi
   sed $OPTS  "s/${kafka_ip}/${ipforreplace}/g" $file
 done
}
