#!/bin/bash +x

#set -e

CHANNEL_NAME=$1
: ${CHANNEL_NAME:="channel1"}
echo $CHANNEL_NAME

export FABRIC_ROOT=$PWD/../..
export FABRIC_CFG_PATH=$PWD
echo

OS_ARCH=$(echo "$(uname -s)-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')

which cryptogen 
which configtxgen

## Using docker-compose template replace private key file names with constants
function replacePrivateKey () {
	ARCH=`uname -s | grep Darwin`
	if [ "$ARCH" == "Darwin" ]; then
		OPTS="-it"
	else
		OPTS="-i"
	fi

	cp docker-compose-e2e-template.yaml docker-compose-e2e.yaml

    CURRENT_DIR=$PWD
    cd crypto-config/peerOrganizations/org1.example.com/ca/
    PRIV_KEY=$(ls *_sk)
    cd $CURRENT_DIR
    sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

    cd crypto-config/peerOrganizations/org2.example.com/ca/
    PRIV_KEY=$(ls *_sk)
    cd $CURRENT_DIR
    sed $OPTS "s/CA2_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

}

## Generates Org certs using cryptogen tool
function generateCerts (){
	CRYPTOGEN=$FABRIC_ROOT/release/$OS_ARCH/bin/cryptogen
	CRYPTOGEN=cryptogen

	if [ -f "$CRYPTOGEN" ]; then
            echo "Using cryptogen -> $CRYPTOGEN"
	else
	    echo "Building cryptogen"
	    make -C $FABRIC_ROOT release-all
	fi

	echo
	echo "##########################################################"
	echo "##### Generate certificates using cryptogen tool #########"
	echo "##########################################################"
	$CRYPTOGEN generate --config=./crypto-config.yaml
	echo
}

## Generate orderer genesis block , channel configuration transaction and anchor peer update transactions
function generateChannelArtifacts() {

	CONFIGTXGEN=$FABRIC_ROOT/release/$OS_ARCH/bin/configtxgen
	CONFIGTXGEN=configtxgen
	if [ -f "$CONFIGTXGEN" ]; then
            echo "Using configtxgen -> $CONFIGTXGEN"
	else
	    echo "Building configtxgen"
	    make -C $FABRIC_ROOT release-all
	fi

	echo "##########################################################"
	echo "#########  Generating Orderer Genesis block ##############"
	echo "##########################################################"
	$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./channel-artifacts/orderer.genesis.block

	echo
	echo "#################################################################"
	echo "### Generating channel configuration transaction 'channel.tx' ###"
	echo "#################################################################"
	$CONFIGTXGEN -profile Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME

	echo
	echo "#################################################################"
	echo "### Generating anchor peer update for Org1MSP/Org2MSP/Org3MSP ###"
	echo "#################################################################"
	for orgMsp in Org1MSP Org2MSP Org3MSP Org4MSP; do
	    $CONFIGTXGEN -profile Channel -outputAnchorPeersUpdate ./channel-artifacts/${orgMsp}anchors.tx -channelID $CHANNEL_NAME -asOrg $orgMsp
    done
	echo
	echo
}


sudo rm -rf crypto-config 
sudo rm -rf channel-artifacts/*
generateCerts
#replacePrivateKey
generateChannelArtifacts

