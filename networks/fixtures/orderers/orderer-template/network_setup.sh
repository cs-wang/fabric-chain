#!/bin/bash

function usage () {
	echo
	echo "======================================================================================================"
	echo "Usage: "
	echo "      network_setup.sh -s [-n <container_name>] <up|down|retstart>"
	echo
	echo "      ./network_setup.sh -s restart"
	echo
	echo "		-s       Enable TLS"
	echo "		-n       The Network Name"
	echo "		up       Launch the network and start the test"
	echo "		down     teardown the network and the test"
	echo "		restart  Restart the network and start the test"
	echo "======================================================================================================"
	echo
}

##process all the options
while getopts "s:h:n:" opt; do
  case "${opt}" in
    s)
      SECURITY="y" #Enable TLS
      ;;
    n)
      CONTAINER_NAME=$OPTARG
      ;;
    h)
      usage
      exit 1
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      exit 1
      ;;
  esac
done

## this is to read the argument up/down/restart
shift $((OPTIND-1))
UP_DOWN="$@"

##Set Defaults
: ${SECURITY:="n"}
: ${CONTAINER_NAME:="testnet"}
: ${COMPOSE_FILE:="docker-compose.yaml"}
: ${UP_DOWN:="restart"}

function clearContainers () {
        #CONTAINER_IDS=$(docker ps -aq)
        CONTAINER_IDS=$(docker ps -a | grep $CONTAINER_NAME | awk '{print $1}')
        if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" = " " ]; then
                echo "---- No containers available for deletion ----"
        else
                docker rm -f $CONTAINER_IDS
        fi
}

function removeUnwantedImages() {
        DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
        if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" = " " ]; then
                echo "---- No images available for deletion ----"
        else
                docker rmi -f $DOCKER_IMAGE_IDS
        fi
}

function networkUp () {
    if [ "$SECURITY" == "y" -o "$SECURITY" == "Y" ]; then
        SECURITY=true
    else
        SECURITY=false
    fi
    ENABLE_TLS=$SECURITY docker-compose -f $COMPOSE_FILE up -d 2>&1

    if [ $? -ne 0 ]; then
	    echo "ERROR !!!! "
	    exit 1
    fi
}

function networkDown () {
    docker-compose -f $COMPOSE_FILE down

    #Cleanup the chaincode containers
    clearContainers

    #Cleanup images
    removeUnwantedImages

}

#Create the network using docker compose
if [ "${UP_DOWN}" == "up" ]; then
	networkUp
elif [ "${UP_DOWN}" == "down" ]; then ## Clear the network
	networkDown
elif [ "${UP_DOWN}" == "restart" ]; then ## Restart the network
	networkDown
	networkUp
else
	usage
	exit 1
fi
