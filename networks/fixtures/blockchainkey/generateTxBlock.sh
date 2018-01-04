#!/bin/bash +x

#set -e

. ./../configall.sh

CHANNEL_NAME=$1
: ${CHANNEL_NAME:="channel1"}
echo $CHANNEL_NAME

## Generate orderer genesis block , channel configuration transaction and anchor peer update transactions
function generateChannelArtifacts() {

	CONFIGTXGEN=configtxgen
	tool=$(which $CONFIGTXGEN)
	if [ -f "$tool" ]; then
        echo "Using configtxgen -> $tool"

	else
        echo "CAN NOT FIND TOOL $CONFIGTXGEN"
        exit 0
	fi

    echo
	echo "##########################################################"
	echo "#########  Generating Orderer Genesis block ##############"
	echo "##########################################################"
	$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./channel-artifacts/orderer.genesis.block -channelID $network_name

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

sudo rm -rf channel-artifacts/*
generateChannelArtifacts

