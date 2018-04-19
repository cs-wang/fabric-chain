#!/bin/bash
#set -x

CHANNEL_NAME="$2"
: ${CHANNEL_NAME:="channel1"}
CORE_PEER_TLS_ENABLED=false

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "Channel Name: "$CHANNEL_NAME

verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
                echo "================== ERROR !!! FAILED to execute reconfig channel Scenario =================="
		echo
   		exit 1
	fi
}

setGlobals () {

       if [ $1 -eq 0 ] ; then
               CORE_PEER_LOCALMSPID="Org1MSP"
               CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
               CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
               CORE_PEER_ADDRESS=peer0.org1.example.com:7051
       elif [ $1 -eq 1 ]; then
	       CORE_PEER_LOCALMSPID="Org2MSP"
               CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
               CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
               CORE_PEER_ADDRESS=peer0.org2.example.com:7051
       elif [ $1 -eq 2 ]; then
               CORE_PEER_LOCALMSPID="Org3MSP"
               CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
               CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
               CORE_PEER_ADDRESS=peer0.org3.example.com:7051
       elif [ $1 -eq 3 ]; then
               CORE_PEER_LOCALMSPID="Org4MSP"
               CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org4.example.com/peers/peer0.org4.example.com/tls/ca.crt
               CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org4.example.com/users/Admin@org4.example.com/msp
               CORE_PEER_ADDRESS=peer0.org4.example.com:7051
       else
		echo "WARNING...NO SUCH ORG..."
       fi

       env |grep CORE
}


fetchChannel() {
	setGlobals 1

	peer channel fetch config config_block.pb -o orderer0.example.com:7050 -c $CHANNEL_NAME >& log.txt
        res=$?
        cat log.txt
        verifyResult $res "Channel fetch failed"
        echo "=====================fetch Channel successfully ===================== "
        echo

        mv config_block.pb ./scripts
}

updateChannel() {
    #    mv ./scripts/config_update_as_envelope.pb ./

        #setGlobals 0
	#peer channel signconfigtx -f config_update_as_envelope2.pb -c $CHANNEL_NAME --cafile crypto/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
	#setGlobals 1
	#peer channel signconfigtx -f config_update_as_envelope2.pb -c $CHANNEL_NAME --cafile crypto/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
	#setGlobals 2
	#peer channel signconfigtx -f config_update_as_envelope2.pb -c $CHANNEL_NAME --cafile crypto/peerOrganizations/org3.example.com/ca/ca.org3.example.com-cert.pem

	CHANNEL_NAME=example
	CORE_PEER_LOCALMSPID="OrdererMSP"
        CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/users/Admin\@example.com/msp
	peer channel signconfigtx -f config_update_as_envelope.pb -c $CHANNEL_NAME --cafile crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/cacerts/ca.example.com-cert.pem
        
	peer channel update -f config_update_as_envelope.pb -c $CHANNEL_NAME -o orderer0.example.com:7050 >& log.txt

        res=$?
        cat log.txt
        verifyResult $res "Channel update failed"
        echo "=====================update Channel successfully ===================== "
        echo

}


if [ "$1" = "1" ]; then
 fetchChannel
elif [ "$1" = "2" ]; then
 updateChannel
else
 echo "error..."
fi

exit
