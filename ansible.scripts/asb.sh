if [ "$1" = "" ] || [ "$1" = "h" ]; then
 echo "please run this script with the first parameter below:"
 echo "	<hosts>: copy the hosts file to /etc/ansible/hosts"
 echo "	<loadimages>: loadimages on orgs node"
 echo "	<dp> <zk>|<orgs>|<key>: deploy zookeeper kafka orderer ca peer couchdb"
 echo "	<cldp>: clear the deployed file"
fi

if [ "$1" = "hosts" ]; then
 ansible_path=/etc/ansible
 if [ ! -d "$ansible_path" ]; then
  echo "create path $ansible_path"
  sudo mkdir $ansible_path
 fi
 sudo cp ./hosts /etc/ansible/hosts
fi

zookeeper0="101.89.66.163"
zookeeper1="101.89.66.164"
zookeeper2="101.89.66.231"
kafka0="101.89.66.163"
kafka1="101.89.66.164"
kafka2="101.89.66.231"

orderer0="101.89.66.163"
orderer1="101.89.66.164"
orderer2="101.89.66.231"

ca1_org1="101.89.66.183"
ca2_org2="101.89.66.184"

peer0_org1="101.89.66.183"
peer1_org1="101.89.66.184"
peer0_org2="101.89.66.184"
peer1_org2="101.89.66.183"

. ../networks/fixtures/configall.sh

if [ "$1" = "loadimages" ]; then
# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-baseos-x86_64-0.3.1.tar"
 ansible orgs -m shell -a "docker tag 4b0cab202084 hyperledger/fabric-baseos:x86_64-0.3.1"

# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-ca-x86_64-1.0.0.tar"
 ansible orgs -m shell -a "docker tag a15c59ecda5b hyperledger/fabric-ca:x86_64-1.0.0"

# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-ccenv-x86_64-1.0.0.tar"
 ansible orgs -m shell -a "docker tag 7182c260a5ca hyperledger/fabric-ccenv:x86_64-1.0.0"

# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-couchdb-x86_64-1.0.0.tar"
 ansible orgs -m shell -a "docker tag 2fbdbf3ab945 hyperledger/fabric-couchdb:x86_64-1.0.0"

# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-orderer-x86_64-1.0.0.tar"
 ansible orgs -m shell -a "docker tag e317ca5638ba hyperledger/fabric-orderer:x86_64-1.0.0"

# ansible orgs -m shell -a "cd /var/blockchain/docker-images && docker load -i fabric-peer-x86_64-1.0.0.tar"
 ansible orgs -m shell -a "docker tag 6830dcd7b9b5 hyperledger/fabric-peer:x86_64-1.0.0"

fi

if [ "$1" = "dp" ]; then

if [ "$2" = "zk" ]; then
 ansible $zookeeper0 -m copy -a "src=../networks/fixtures/zkfks/zookeeper0 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
 ansible $zookeeper0 -m copy -a "src=../networks/fixtures/zkfks/kafka0 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
 ansible $zookeeper1 -m copy -a "src=../networks/fixtures/zkfks/zookeeper1 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
 ansible $zookeeper1 -m copy -a "src=../networks/fixtures/zkfks/kafka1 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
 ansible $zookeeper2 -m copy -a "src=../networks/fixtures/zkfks/zookeeper2 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
 ansible $zookeeper2 -m copy -a "src=../networks/fixtures/zkfks/kafka2 dest=/var/blockchain/$network_name/networks/fixtures/zkfks/ owner=root mode=0755"
fi

if [ "$2" = "od" ]; then
 ansible $orderer0 -m copy -a "src=../networks/fixtures/orderers/orderer0 dest=/var/blockchain/$network_name/networks/fixtures/orderers/ owner=root mode=0755"
 ansible $orderer1 -m copy -a "src=../networks/fixtures/orderers/orderer1 dest=/var/blockchain/$network_name/networks/fixtures/orderers/ owner=root mode=0755"
 ansible $orderer2 -m copy -a "src=../networks/fixtures/orderers/orderer2 dest=/var/blockchain/$network_name/networks/fixtures/orderers/ owner=root mode=0755"
fi
if [ "$2" = "orgs" ]; then
  set -x
  ansible $ca1_org1 -m copy -a "src=../networks/fixtures/orgs/org1/ca1 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org1 owner=root mode=0755"
  ansible $peer0_org1 -m copy -a "src=../networks/fixtures/orgs/org1/peer0 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org1 owner=root mode=0755"
  ansible $peer1_org1 -m copy -a "src=../networks/fixtures/orgs/org1/peer1 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org1 owner=root mode=0755"
  ansible $ca2_org2 -m copy -a "src=../networks/fixtures/orgs/org2/ca2 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org2 owner=root mode=0755"
  ansible $peer0_org2 -m copy -a "src=../networks/fixtures/orgs/org2/peer0 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org2 owner=root mode=0755"
  ansible $peer1_org2 -m copy -a "src=../networks/fixtures/orgs/org2/peer1 dest=/var/blockchain/$network_name/networks/fixtures/orgs/org2 owner=root mode=0755"

fi

if [ "$2" = "key" ]; then
set -x
 CURRENT_DIR=$PWD
 cd ../networks/fixtures/blockchainkey/
 tar -czf crypto-config.tar.gz crypto-config/
 tar -czf channel-artifacts.tar.gz channel-artifacts/
 cd $CURRENT_DIR
 ansible orderers -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/blockchainkey/ && rm -rf *"
 ansible orderers -m copy -a "src=../networks/fixtures/blockchainkey/crypto-config.tar.gz dest=/var/blockchain/$network_name/networks/fixtures/blockchainkey/ owner=root mode=0755"
 ansible orderers -m copy -a "src=../networks/fixtures/blockchainkey/channel-artifacts.tar.gz dest=/var/blockchain/$network_name/networks/fixtures/blockchainkey/ owner=root mode=0755"
 ansible orderers -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/blockchainkey/ && tar -zxf crypto-config.tar.gz && rm -f crypto-config.tar.gz"
 ansible orderers -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/blockchainkey/ && tar -zxf channel-artifacts.tar.gz && rm -f channel-artifacts.tar.gz"
 
 ansible orgs -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/blockchainkey/ && rm -rf *"
 ansible orgs -m copy -a "src=../networks/fixtures/blockchainkey/crypto-config.tar.gz dest=/var/blockchain/$network_name/networks/fixtures/blockchainkey/ owner=root mode=0755"
 ansible orgs -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/blockchainkey/ && tar -zxf crypto-config.tar.gz && rm -f crypto-config.tar.gz"
 cd ../networks/fixtures/blockchainkey/
 rm -f crypto-config.tar.gz
 rm -f channel-artifacts.tar.gz
 cd $CURRENT_DIR
fi

fi

if [ "$1" = "cldp" ]; then
 ansible zkfks -m shell -a "cd /var/blockchain && rm -rf $network_name"
 ansible orgs -m shell -a "cd /var/blockchain && rm -rf $network_name"
fi

if [ "$1" = "up" ]; then
 if [ "$2" = "zk" ] || [ "$2" = "" ]; then
 ansible $zookeeper0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper0 && ./network_setup.sh up"
 ansible $zookeeper1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper1 && ./network_setup.sh up"
 ansible $zookeeper2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper2 && ./network_setup.sh up"
 ansible $kafka0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka0 && ./network_setup.sh up"
 ansible $kafka1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka1 && ./network_setup.sh up"
 ansible $kafka2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka2 && ./network_setup.sh up"
 fi
 
 if [ "$2" = "od" ] || [ "$2" = "" ]; then
  ansible $orderer0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer0 && ./network_setup.sh -s up"
  ansible $orderer1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer1 && ./network_setup.sh -s up"
  ansible $orderer2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer2 && ./network_setup.sh -s up"
 fi
 
 if [ "$2" = "orgs" ] || [ "$2" = "" ]; then
  set -x
  ansible orgs -m shell -a "chmod -R 777 /paic_data"

  ansible $ca1_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/ca1 && ./network_setup.sh -s up"
  ansible $peer0_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/peer0 && ./network_setup.sh -s up"
  ansible $peer1_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/peer1 && ./network_setup.sh -s up"
  ansible $ca2_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/ca2 && ./network_setup.sh -s up"
  ansible $peer0_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/peer0 && ./network_setup.sh -s up"
  ansible $peer1_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/peer1 && ./network_setup.sh -s up"

 fi
fi

if [ "$1" = "down" ]; then
 if [ "$2" = "zk" ] || [ "$2" = "" ]; then
 ansible $kafka0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka0 && ./network_setup.sh down"
 ansible $kafka1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka1 && ./network_setup.sh down"
 ansible $kafka2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/kafka2 && ./network_setup.sh down"
 ansible $zookeeper0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper0 && ./network_setup.sh down"
 ansible $zookeeper1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper1 && ./network_setup.sh down"
 ansible $zookeeper2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/zkfks/zookeeper2 && ./network_setup.sh down"
 fi
 
 if [ "$2" = "od" ] || [ "$2" = "" ]; then
  ansible $orderer0 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer0 && ./network_setup.sh -n orderer0.$network_name.com down"
  ansible $orderer1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer1 && ./network_setup.sh -n orderer1.$network_name.com down"
  ansible $orderer2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers/orderer2 && ./network_setup.sh -n orderer2.$network_name.com down"
 fi
 
 if [ "$2" = "orgs" ] || [ "$2" = "" ]; then

  ansible $ca1_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/ca1 && ./network_setup.sh -n ca1 down"
  ansible $peer0_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/peer0 && ./network_setup.sh -n peer0.org1.$network_name.com down"
  ansible $peer1_org1 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org1/peer1 && ./network_setup.sh -n peer1.org1.$network_name.com down"
  ansible $ca2_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/ca2 && ./network_setup.sh -n ca2 down"
  ansible $peer0_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/peer0 && ./network_setup.sh -n peer0.org2.$network_name.com down"
  ansible $peer1_org2 -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/org2/peer1 && ./network_setup.sh -n peer1.org2.$network_name.com down"

 fi
fi

if [ "$1" = "clda" ]; then
 if [ "$2" = "zk" ] || [ "$2" = "" ]; then
 ansible zkfks -m shell -a "rm -rf /paic_data/*"
 fi
 if [ "$2" = "orgs" ] || [ "$2" = "" ]; then
 ansible orgs -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orgs/data && for couchdb in \$(find -name *couchdb*); do rm -rf \${couchdb}/* ; done && rm -rf \$(find -name *fabric*)"
 fi
 if [ "$2" = "od" ] || [ "$2" = "" ]; then
 ansible orderers -m shell -a "cd /var/blockchain/$network_name/networks/fixtures/orderers && rm -rf data/"
 fi
fi

if [ "$1" = "compose" ]; then
 ansible zkfks -m copy -a "src=~/docker-compose dest=/var/blockchain/ owner=root mode=0755"
 ansible orgs -m copy -a "src=~/docker-compose dest=/var/blockchain/ owner=root mode=0755"
 ansible zkfks -m shell -a "cd /var/blockchain && mv docker-compose /usr/local/bin/"
 ansible orgs -m shell -a "cd /var/blockchain && mv docker-compose /usr/local/bin/"
fi


