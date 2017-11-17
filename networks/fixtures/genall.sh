#set -x
. configall.sh
#set -x
function gentmpl() {
 
 cp $orderer_template_path/$src_template $orderer_template_path/$dest_template
 cp $peer_template_path/$src_template $peer_template_path/$dest_template
 
 cp $zookeeper_template_path/$src_template $zookeeper_template_path/$dest_template
 cp $kafka_template_path/$src_template $kafka_template_path/$dest_template

}

function initop() {

 for template in $orderer_template_path/$dest_template $peer_template_path/$dest_template; do
  if [ "$net_mode" = "local" ]; then
    echo "local mode test"
  fi
  echo "    "$extra_host >> $template

  for orderer in ${ORDERERS[@]}; do
   #orderer_host=`eval echo '$'"${orderer}"_"host"`
   orderer_domain="$orderer.$network_name.com"
   orderer_ip="${orderer}_${host_end}"
   orderer_host="${orderer_domain}:${orderer_ip}"
   echo "$blank"- $orderer_host   >> $template
  done
  for org in ${ORGS[@]}; do
   for peer in ${PEERS[@]}; do
    #peer_org_host=`eval echo '$'"${peer}"_"${org}"_"host"`
    peer_domain="${peer}.${org}.${network_name}.com"
    peer_ip="${peer}_${org}_${host_end}"
    peer_org_host="${peer_domain}:${peer_ip}"
    echo "$blank"- $peer_org_host >> $template
   done
  done
  for kafka in ${KAFKAS[@]}; do
   #kafka_host=`eval echo '$'"${kafka}"_"host"`
   kafka_domain="${kafka}"
   kafka_ip="${kafka}_${host_end}"
   kafka_host="${kafka_domain}:${kafka_ip}"
   echo "$blank"- $kafka_host >> $template
  done

 done
}

function initzk() {

 for template in $zookeeper_template_path/$dest_template $kafka_template_path/$dest_template; do
  echo "    "$extra_host >> $template
  for zookeeper in ${ZOOKEEPERS[@]}; do
   #zookeeper_host=`eval echo '$'"${zookeeper}"_"host"`
   zookeeper_domain="${zookeeper}"
   zookeeper_ip="${zookeeper}_${host_end}"
   zookeeper_host="${zookeeper_domain}:${zookeeper_ip}"
   echo "$blank"- $zookeeper_host   >> $template
  done

  for kafka in ${KAFKAS[@]}; do
   kafka_host=`eval echo '$'"${kafka}"_"host"`
   kafka_domain="${kafka}"
   kafka_ip="${kafka}_${host_end}"
   kafka_host="${kafka_domain}:${kafka_ip}"
   echo "$blank"- $kafka_host >> $template
  done

 done

}

function inithost() {

 initop
 
 initzk

}

function main() {
 if [ "$1" = "" ]; then
  cd orderers
  ./generateorderers.sh
  cd ../
  cd orgs/
  ./generateorgs.sh
  cd ../
  cd zkfks/
  ./generatezkfks.sh
  cd ../
 else 
  cd orderers
  ./generateorderers.sh clean
  cd ../
  cd orgs/
  ./generateorgs.sh clean
  cd ../
  cd zkfks/
  ./generatezkfks.sh clean
  cd ../
 fi
}

gentmpl
inithost
main $1
