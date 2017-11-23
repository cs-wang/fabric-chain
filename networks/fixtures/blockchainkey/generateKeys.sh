#!/bin/bash +x

#set -e

CHANNEL_NAME=$1
: ${CHANNEL_NAME:="channel1"}
echo $CHANNEL_NAME

## Generates Org certs using cryptogen tool
function generateCerts (){
	CRYPTOGEN=cryptogen
    tool=$(which $CRYPTOGEN)
	if [ -f "$tool" ]; then
            echo "Using cryptogen -> $tool"
	else
        echo "CAN NOT FIND TOOL $CRYPTOGEN"
        exit 0
	fi

	echo
	echo "##########################################################"
	echo "##### Generate certificates using cryptogen tool #########"
	echo "##########################################################"
	$CRYPTOGEN generate --config=./crypto-config.yaml
	echo
}

sudo rm -rf crypto-config
generateCerts