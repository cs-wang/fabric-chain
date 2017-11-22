/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	//"testing"
	"time"

	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	client "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/admin"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"os"
	"sync"
)

var org1 = "peerorg1"
var org2 = "peerorg2"
var org3 = "peerorg3"
var org4 = "peerorg4"
// Client
var orgTestClient fab.FabricClient

// Channel
var orgTestChannel1 fab.Channel
//var orgTestChannel2 fab.Channel

// Orderers
var orgTestOrderer fab.Orderer

// Peers
var org1TestPeer0 fab.Peer
//var org1TestPeer1 fab.Peer

var org2TestPeer0 fab.Peer
//var org2TestPeer1 fab.Peer

var org3TestPeer0 fab.Peer
//var org3TestPeer1 fab.Peer
var org4TestPeer0 fab.Peer

// EventHubs
var org1peer0EventHub fab.EventHub
//var org1peer1EventHub fab.EventHub

var org2peer0EventHub fab.EventHub
//var org2peer1EventHub fab.EventHub

var org3peer0EventHub fab.EventHub
//var org3peer1EventHub fab.EventHub
var org4peer0EventHub fab.EventHub

// Users
var org1AdminUser ca.User
var org2AdminUser ca.User
var org3AdminUser ca.User
var org4AdminUser ca.User
var ordererAdminUser ca.User
var org1User ca.User
var org2User ca.User
var org3User ca.User
var org4User ca.User

var gchainCodeID = "marbles02"
var endorseAdminUser = &org1AdminUser
var endorsePeer = &org1TestPeer0

var successNumLock sync.RWMutex
var invokeSuccessNum int64 = 0

// initializeFabricClient initializes fabric-sdk-go
func initializeFabricClient() {
	// Initialize configuration
	configImpl, err := config.InitConfig("./../config/config_netapp.yaml")
	failTestIfError(err, "configImpl, err := config.InitConfig")

	// Instantiate client
	orgTestClient = client.NewClient(configImpl)

	// Initialize crypto suite
	err = factory.InitFactories(configImpl.CSPConfig())
	failTestIfError(err, "factory.InitFactories")
	cryptoSuite := factory.GetDefault()
	orgTestClient.SetCryptoSuite(cryptoSuite)

	StateStore, err := fabapi.NewDefaultStateStore(
		fabapi.StateStoreOpts{
			Path: "/tmp/enroll_user",
		},
		orgTestClient.Config())
	failTestIfError(err, "fabapi.NewDefaultStateStore")
	orgTestClient.SetStateStore(StateStore)
}

func createTestChannel() {
	var err error

	orgTestChannel1, err = channel.NewChannel("channel2", orgTestClient)
	failTestIfError(err, "orgTestChannel1, err = channel.NewChannel")
	//orgTestChannel2, err = channel.NewChannel("channel2", orgTestClient)
	//failTestIfError(err, "orgTestChannel2, err = channel.NewChannel")

	err = orgTestChannel1.AddPeer(org1TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org1TestPeer0)")
	err = orgTestChannel1.AddPeer(org2TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org2TestPeer0)")
	err = orgTestChannel1.AddPeer(org3TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org3TestPeer0)")
	err = orgTestChannel1.AddPeer(org4TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org4TestPeer0)")
	err = orgTestChannel1.AddOrderer(orgTestOrderer)
	failTestIfError(err, "orgTestChannel1.AddOrderer(orgTestOrderer)")
	err = orgTestChannel1.SetPrimaryPeer(org1TestPeer0)
	failTestIfError(err, "orgTestChannel1.SetPrimaryPeer(org1TestPeer0)")

	//err = admin.CreateOrUpdateChannel(orgTestClient, ordererAdminUser, org1AdminUser,
	//	orgTestChannel1, "./../fixtures/channel/config/channel1.tx")
	err = admin.CreateOrUpdateChannel(orgTestClient, ordererAdminUser, org1AdminUser,
		orgTestChannel1, "./../../fixtures/blockchainkey/channel-artifacts/mychannel.tx")
	failTestIfError(err, "CreateOrUpdateChannel1")
	// Allow orderer to process channel creation
	time.Sleep(time.Millisecond * 5000)

	err = orgTestChannel1.Initialize(nil)
	failTestIfError(err, "orgTestChannel1.Initialize")
}

func loadTestChannel() {
	var err error

	orgTestChannel1, err = channel.NewChannel("channel4", orgTestClient)
	failTestIfError(err, "orgTestChannel1, err = channel.NewChannel")

	err = orgTestChannel1.AddPeer(org1TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org1TestPeer0)")
	err = orgTestChannel1.AddPeer(org2TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org2TestPeer0)")
	err = orgTestChannel1.AddPeer(org3TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org3TestPeer0)")
	err = orgTestChannel1.AddPeer(org4TestPeer0)
	failTestIfError(err, "orgTestChannel1.AddPeer(org4TestPeer0)")
	err = orgTestChannel1.AddOrderer(orgTestOrderer)
	failTestIfError(err, "orgTestChannel1.AddOrderer(orgTestOrderer)")
	err = orgTestChannel1.SetPrimaryPeer(org1TestPeer0)
	failTestIfError(err, "orgTestChannel1.SetPrimaryPeer(org1TestPeer0)")
}

func joinTestChannel() {
	// Get org1peer0 to join channel
	orgTestChannel1.RemovePeer(org1TestPeer0)
	orgTestChannel1.RemovePeer(org2TestPeer0)
	orgTestChannel1.RemovePeer(org3TestPeer0)
	orgTestChannel1.RemovePeer(org4TestPeer0)

	orgTestChannel1.AddPeer(org1TestPeer0)
	orgTestChannel1.SetPrimaryPeer(org1TestPeer0)
	err := admin.JoinChannel(orgTestClient, org1AdminUser, orgTestChannel1)
	failTestIfError(err, "org1 JoinChannel1")

	// Get org2peer0 to join channel
	orgTestChannel1.RemovePeer(org1TestPeer0)
	orgTestChannel1.AddPeer(org2TestPeer0)
	orgTestChannel1.SetPrimaryPeer(org2TestPeer0)
	err = admin.JoinChannel(orgTestClient, org2AdminUser, orgTestChannel1)
	failTestIfError(err, "org2 JoinChanne1")

	// Get org3peer0 to join channel
	orgTestChannel1.RemovePeer(org2TestPeer0)
	orgTestChannel1.AddPeer(org3TestPeer0)
	orgTestChannel1.SetPrimaryPeer(org3TestPeer0)
	err = admin.JoinChannel(orgTestClient, org3AdminUser, orgTestChannel1)
	failTestIfError(err, "org3 JoinChanne1")

	// Get org4peer0 to join channel
	orgTestChannel1.RemovePeer(org3TestPeer0)
	orgTestChannel1.AddPeer(org4TestPeer0)
	orgTestChannel1.SetPrimaryPeer(org4TestPeer0)
	err = admin.JoinChannel(orgTestClient, org4AdminUser, orgTestChannel1)
	failTestIfError(err, "org4 JoinChanne1")

	orgTestChannel1.RemovePeer(org4TestPeer0)
}

func installAndInstantiate() {

	orgTestChannel1.AddPeer(org1TestPeer0)
	orgTestChannel1.AddPeer(org2TestPeer0)
	orgTestChannel1.AddPeer(org3TestPeer0)
	orgTestChannel1.AddPeer(org4TestPeer0)
	var err error
	fmt.Println("=================1=================")
	orgTestClient.SetUserContext(org1AdminUser)
	err =admin.SendInstallCC(orgTestClient, gchainCodeID,
		"example", "0", nil, []fab.Peer{org1TestPeer0}, "./../chaincode")
	failTestIfError(err, "SendInstallCC org1TestPeer0")

	fmt.Println("=================2=================")
	orgTestClient.SetUserContext(org2AdminUser)
	err = admin.SendInstallCC(orgTestClient, gchainCodeID,
		"example", "0", nil, []fab.Peer{org2TestPeer0}, "./../chaincode")
	failTestIfError(err, "SendInstallCC org2TestPeer0")

	fmt.Println("=================3=================")
	orgTestClient.SetUserContext(org3AdminUser)
	err = admin.SendInstallCC(orgTestClient, gchainCodeID,
		"example", "0", nil, []fab.Peer{org3TestPeer0}, "./../chaincode")
	failTestIfError(err, "SendInstallCC org3TestPeer0")

	fmt.Println("=================4=================")
	orgTestClient.SetUserContext(org4AdminUser)
	err = admin.SendInstallCC(orgTestClient, gchainCodeID,
		"example", "0", nil, []fab.Peer{org4TestPeer0}, "./../chaincode")
	failTestIfError(err, "SendInstallCC org3TestPeer0")

	fmt.Println("111111111111111111111111111111111111")

	orgTestClient.SetUserContext(*endorseAdminUser)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	err = admin.SendInstantiateCC(orgTestChannel1, gchainCodeID,
		generateInitArgs(), "example", "0", []apitxn.ProposalProcessor{ org1TestPeer0, org2TestPeer0, org3TestPeer0, org4TestPeer0 }, org1peer0EventHub)
	failTestIfError(err, "SendInstantiateCC orgTestChannel1")
	fmt.Println("222222222222222222222222222222222222")
}

func loadOrderer() {
	ordererConfig, err := orgTestClient.Config().RandomOrdererConfig()
	failTestIfError(err, "ordererConfig, err := orgTestClient.Config().RandomOrdererConfig()")

	orgTestOrderer, err = orderer.NewOrderer(fmt.Sprintf("%s:%d", ordererConfig.Host,
		ordererConfig.Port), ordererConfig.TLS.Certificate,
		ordererConfig.TLS.ServerHostOverride, orgTestClient.Config())
	failTestIfError(err, "orderer.NewOrderer")
}

func loadOrgPeers() {
	org1Peers, err := orgTestClient.Config().PeersConfig(org1)
	failTestIfError(err, "org1Peers, err := orgTestClient.Config")

	org2Peers, err := orgTestClient.Config().PeersConfig(org2)
	failTestIfError(err, "org2Peers, err := orgTestClient.Config")

	org3Peers, err := orgTestClient.Config().PeersConfig(org3)
	failTestIfError(err, "org3Peers, err := orgTestClient.Config")

	org4Peers, err := orgTestClient.Config().PeersConfig(org4)
	failTestIfError(err, "org4Peers, err := orgTestClient.Config")

	org1TestPeer0, err = peer.NewPeerTLSFromCert(fmt.Sprintf("%s:%d", org1Peers[0].Host,
		org1Peers[0].Port), org1Peers[0].TLS.Certificate,
		org1Peers[0].TLS.ServerHostOverride, orgTestClient.Config())
	failTestIfError(err, "org1TestPeer0, err = peer.NewPeerTLSFromCert")

	org2TestPeer0, err = peer.NewPeerTLSFromCert(fmt.Sprintf("%s:%d", org2Peers[0].Host,
		org2Peers[0].Port), org2Peers[0].TLS.Certificate,
		org2Peers[0].TLS.ServerHostOverride, orgTestClient.Config())
	failTestIfError(err, "org2TestPeer0, err = peer.NewPeerTLSFromCert")

	org3TestPeer0, err = peer.NewPeerTLSFromCert(fmt.Sprintf("%s:%d", org3Peers[0].Host,
		org3Peers[0].Port), org3Peers[0].TLS.Certificate,
		org3Peers[0].TLS.ServerHostOverride, orgTestClient.Config())
	failTestIfError(err, "org3TestPeer0, err = peer.NewPeerTLSFromCert")

	org4TestPeer0, err = peer.NewPeerTLSFromCert(fmt.Sprintf("%s:%d", org4Peers[0].Host,
		org4Peers[0].Port), org4Peers[0].TLS.Certificate,
		org4Peers[0].TLS.ServerHostOverride, orgTestClient.Config())
	failTestIfError(err, "org4TestPeer0, err = peer.NewPeerTLSFromCert")


	org1peer0EventHub, err = events.NewEventHub(orgTestClient)
	failTestIfError(err, "org1peer0EventHub, err = events.NewEventHub")

	org1peer0EventHub.SetPeerAddr(fmt.Sprintf("%s:%d", org1Peers[0].EventHost,
		org1Peers[0].EventPort), org1Peers[0].TLS.Certificate,
		org1Peers[0].TLS.ServerHostOverride)

	orgTestClient.SetUserContext(org1User)
	err = org1peer0EventHub.Connect()
	failTestIfError(err, "org1peer0EventHub.Connect")

	org2peer0EventHub, err = events.NewEventHub(orgTestClient)
	failTestIfError(err, "org2peer0EventHub, err = events.NewEventHub")

	org2peer0EventHub.SetPeerAddr(fmt.Sprintf("%s:%d", org2Peers[0].EventHost,
		org2Peers[0].EventPort), org2Peers[0].TLS.Certificate,
		org2Peers[0].TLS.ServerHostOverride)

	orgTestClient.SetUserContext(org2User)
	err = org2peer0EventHub.Connect()

	failTestIfError(err, "org2peer0EventHub.Connect")

	org3peer0EventHub, err = events.NewEventHub(orgTestClient)
	failTestIfError(err, "org3peer0EventHub, err = events.NewEventHub")

	org3peer0EventHub.SetPeerAddr(fmt.Sprintf("%s:%d", org3Peers[0].EventHost,
		org3Peers[0].EventPort), org3Peers[0].TLS.Certificate,
		org3Peers[0].TLS.ServerHostOverride)

	orgTestClient.SetUserContext(org3User)
	err = org3peer0EventHub.Connect()

	failTestIfError(err, "org3peer0EventHub.Connect")

	org4peer0EventHub, err = events.NewEventHub(orgTestClient)
	failTestIfError(err, "org4peer0EventHub, err = events.NewEventHub")

	org4peer0EventHub.SetPeerAddr(fmt.Sprintf("%s:%d", org4Peers[0].EventHost,
		org4Peers[0].EventPort), org4Peers[0].TLS.Certificate,
		org4Peers[0].TLS.ServerHostOverride)

	orgTestClient.SetUserContext(org4User)
	err = org4peer0EventHub.Connect()
	failTestIfError(err, "org4peer0EventHub.Connect")
}

// loadOrgUsers Loads all the users required to perform this test
func loadOrgUsers() {
	var err error

	ordererAdminUser, err = GetOrdererAdmin(orgTestClient, org1)
	failTestIfError(err, "GetOrdererAdmin ordererAdminUser")
	org1AdminUser, err = GetAdmin(orgTestClient, "org1", org1)
	failTestIfError(err, "GetAdmin org1AdminUser")
	org2AdminUser, err = GetAdmin(orgTestClient, "org2", org2)
	failTestIfError(err, "GetAdmin org2AdminUser")
	org3AdminUser, err = GetAdmin(orgTestClient, "org3", org3)
	failTestIfError(err, "GetAdmin org3AdminUser")
	org4AdminUser, err = GetAdmin(orgTestClient, "org4", org4)
	failTestIfError(err, "GetAdmin org4AdminUser")
	org1User, err = GetUser(orgTestClient, "org1", org1)
	failTestIfError(err, "GetUser, org1User")
	org2User, err = GetUser(orgTestClient, "org2", org2)
	failTestIfError(err, "GetUser, org2User")
	org3User, err = GetUser(orgTestClient, "org3", org3)
	failTestIfError(err, "GetUser, org3User")
	org4User, err = GetUser(orgTestClient, "org4", org4)
	failTestIfError(err, "GetUser, org4User")
}

func generateInitArgs() []string {
	var args []string
	args = append(args, "init")
	args = append(args, "a")
	args = append(args, "100")
	args = append(args, "b")
	args = append(args, "200")
	return nil
}

func failTestIfError(err error, msg string) {
	if err != nil {
		fmt.Print("============msg========= : [", msg, "], err : ")
		fmt.Println(err)
		os.Exit(1)
	}
}

func failInvokeTestIfError(err error, msg string) {
	if err != nil {
		fmt.Print("============invoke failed msg=========\n")
		fmt.Println(err)
	} else {
		successNumLock.Lock()
		defer successNumLock.Unlock()
		invokeSuccessNum++
	}

}
