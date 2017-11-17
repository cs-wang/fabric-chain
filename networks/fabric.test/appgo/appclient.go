package main

import (
	"fmt"
	"time"
	"os"
	"path"

	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	"github.com/hyperledger/fabric-sdk-go/api/apifabca"
	"github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/admin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn"
	"strconv"
)

type APPClient struct {
	Client          apifabclient.FabricClient
	ConfigFile      string
	StateStoreOpts  fabapi.StateStoreOpts
	EventHub        apifabclient.EventHub
	ConnectEventHub bool
	Orderer         *orderer.Orderer
}

type Channel struct {
	AppClient     *APPClient
	Channel       apifabclient.Channel
	ChannelId     string
	ChannelOrgIDs []string
	ChannelConfig string
	Organizations []*Organization
	PrimaryOrgID  string
	ChainCode     ChainCode
}

type ChainCode struct {
	ChainCodeId      string
	ChainCodePath    string
	ChainCodeVersion string
}

type Organization struct {
	OrgID      string
	OrgPath    string
	MSPClient  apifabca.FabricCAClient
	AdminUser  apifabca.User
	NormalUser apifabca.User
	orgPeers   []apifabclient.Peer
}

type CurrentContext struct {
	AppClient    *APPClient
	Organization *Organization
	Channel      *Channel
}

// Initialize reads configuration from file and sets up client, channel and event hub
func (appClient *APPClient) InitAPPClient() error {
	appClient.StateStoreOpts = fabapi.StateStoreOpts{
		Path: "/tmp/enroll_user",
	}
	// Initialize default config provider
	if appClient.Client == nil {
		config, err := fabapi.NewDefaultConfig(fabapi.ConfigOpts{}, fabapi.SDKOpts{ConfigFile:appClient.ConfigFile})
		if err != nil {
			return fmt.Errorf("Failed to initialize default config [%s]", err)
		}
		appClient.Client = fabapi.NewSystemClient(config)
	}

	// Initialize default crypto provider
	if appClient.Client.CryptoSuite() == nil {
		cryptoSuite, err := fabapi.NewCryptoSuite(appClient.Client.Config().CSPConfig())
		if err != nil {
			return fmt.Errorf("Failed to initialize default crypto suite [%s]", err)
		}
		appClient.Client.SetCryptoSuite(cryptoSuite)
	}
	// Initialize default state store
	if appClient.Client.StateStore() == nil {
		store, err := fabapi.NewDefaultStateStore(appClient.StateStoreOpts, appClient.Client.Config())
		if err != nil {
			return fmt.Errorf("Failed to initialize default state store [%s]", err)
		}
		appClient.Client.SetStateStore(store)
	}

	ordererConfig, err := appClient.Client.Config().RandomOrdererConfig()
	if err != nil {
		return fmt.Errorf("RandomOrdererConfig() return error: %s", err)
	}

	orderer, err := orderer.NewOrderer(fmt.Sprintf("%s:%d", ordererConfig.Host,
		ordererConfig.Port), ordererConfig.TLS.Certificate,
		ordererConfig.TLS.ServerHostOverride, appClient.Client.Config())
	if err != nil {
		return fmt.Errorf("NewOrderer return error: %v", err)
	}
	appClient.Orderer = orderer

	eventHub, err := events.NewEventHub(appClient.Client)
	if err != nil {
		return fmt.Errorf("Error creating new event hub: %v", err)
	}
	appClient.EventHub = eventHub
	return nil
}

func (context *CurrentContext) SetOrgUser() error {
	client := context.AppClient
	// Initialize default MSP client
	if context.Organization.MSPClient == nil {
		// TODO: Becomes MSP Manager
		mspClient, err := fabapi.NewCAClient(context.Organization.OrgID, client.Client.Config())
		if err != nil {
			return fmt.Errorf("Failed to initialize default client [%s]", err)
		}
		context.Organization.MSPClient = mspClient
	}
	user, err := fabapi.NewUser(client.Client.Config(), context.Organization.MSPClient, "admin", "adminpw", context.Organization.OrgID)
	if err != nil {
		return fmt.Errorf("NewUser returned error: %v", err)
	}
	err = client.Client.SaveUserToStateStore(user, false)
	if err != nil {
		return fmt.Errorf("client.SaveUserToStateStore returned error: %v", err)
	}
	client.Client.SetUserContext(user)

	orgAdmin, err := GetAdmin(client.Client, context.Organization.OrgPath, context.Organization.OrgID)
	if err != nil {
		return fmt.Errorf("Error getting org admin user: %v", err)
	}
	context.Organization.AdminUser = orgAdmin
	orgUser, err := GetUser(client.Client, context.Organization.OrgPath, context.Organization.OrgID)
	if err != nil {
		return fmt.Errorf("Error getting org user: %v", err)
	}
	context.Organization.NormalUser = orgUser
	return nil
}

func (context *CurrentContext) SetPeers() error {
	client := context.AppClient
	peerConfig, err := client.Client.Config().PeersConfig(context.Organization.OrgID)
	if err != nil {
		fmt.Println(err)
	}
	context.Organization.orgPeers = make([]apifabclient.Peer, len(peerConfig))
	for i, peer := range peerConfig {
		tempPeer, err := fabapi.NewPeer(fmt.Sprintf("%s:%d", peer.Host, peer.Port),
			peer.TLS.Certificate, peer.TLS.ServerHostOverride, client.Client.Config())
		if err != nil {
			fmt.Println(err)
		}
		context.Organization.orgPeers[i] = tempPeer
	}
	return nil
}

func (context *CurrentContext) ConnectEventHub() error {
	client := context.AppClient
	//by default client's user context should use regular user, for admin actions, UserContext must be set to AdminUser
	client.Client.SetUserContext(context.Organization.NormalUser)
	err := context.setEventHub()
	if err != nil {
		return err
	}
	if client.ConnectEventHub {
		if err := client.EventHub.Connect(); err != nil {
			return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
		}
	}

	return nil
}

// getEventHub initilizes the event hub
func (context *CurrentContext) setEventHub() error {
	client := context.AppClient
	foundEventHub := false
	peerConfig, err := client.Client.Config().PeersConfig(context.Organization.OrgID)
	if err != nil {
		return fmt.Errorf("Error reading peer config: %v", err)
	}
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("******* EventHub connect to peer (%s:%d) *******\n", p.EventHost, p.EventPort)
			client.EventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return fmt.Errorf("No EventHub configuration found")
	}

	return nil
}

func (channel *Channel) SetChannel() error {
	tempChannel, err := channel.getChannel(channel.ChannelId, channel.ChannelOrgIDs)
	if err != nil {
		return fmt.Errorf("Create channel (%s) failed: %v", channel.ChannelId, err)
	}
	channel.Channel = tempChannel
	return nil
}

// GetChannel initializes and returns a channel based on config
func (channel *Channel) getChannel(channelID string, orgs []string) (apifabclient.Channel, error) {
	client := channel.AppClient
	channelInner, err := client.Client.NewChannel(channelID)
	if err != nil {
		return nil, fmt.Errorf("NewChannel return error: %v", err)
	}
	err = channelInner.AddOrderer(client.Orderer)
	if err != nil {
		return nil, fmt.Errorf("Error adding orderer: %v", err)
	}
	for _, org := range orgs {
		peerConfig, err := client.Client.Config().PeersConfig(org)
		if err != nil {
			return nil, fmt.Errorf("Error reading peer config: %v", err)
		}
		for _, p := range peerConfig {
			endorser, err := fabapi.NewPeer(fmt.Sprintf("%s:%d", p.Host, p.Port),
				p.TLS.Certificate, p.TLS.ServerHostOverride, client.Client.Config())

			if err != nil {
				return nil, fmt.Errorf("NewPeer return error: %v", err)
			}
			err = channelInner.AddPeer(endorser)
			if err != nil {
				return nil, fmt.Errorf("Error adding peer: %v", err)
			}
			if p.Primary {
				channelInner.SetPrimaryPeer(endorser)
				channel.PrimaryOrgID = org
			}
		}
	}
	return channelInner, nil
}

func (channel *Channel) CreateAndJoinChannel() error {
	client := channel.AppClient
	// Check if primary peer has joined channel
	primaryOrg := getPrimaryPeerOrg(channel)
	//alreadyJoined, err := HasPrimaryPeerJoinedChannel(client.Client, primaryOrg.AdminUser, channel.Channel)
	//if err != nil {
	//	return fmt.Errorf("Error while checking if primary peer has already joined channel: %v", err)
	//}
	logger.Info("CreateAndJoinChannel GetOrdererAdmin")
	ordererAdmin, err := GetOrdererAdmin(client.Client, primaryOrg.OrgID)
	if err != nil {
		return fmt.Errorf("Error getting orderer admin user: %v", err)
	}

	logger.Info("CreateAndJoinChannel CreateOrUpdateChannel")
	// Create, initialize and join channel
	if err = admin.CreateOrUpdateChannel(client.Client, ordererAdmin, primaryOrg.AdminUser, channel.Channel, channel.ChannelConfig); err != nil {
		return fmt.Errorf("CreateChannel returned error: %v", err)
	}
	time.Sleep(time.Second * 3)
	client.Client.SetUserContext(primaryOrg.AdminUser)
	if err = channel.Channel.Initialize(nil); err != nil {
		return fmt.Errorf("Error initializing channel: %v", err)
	}
	//if !alreadyJoined {
	//	logger.Info("alreadyJoined ====",alreadyJoined)
	for i, org := range channel.Organizations {
		logger.Info("channel.Organizations i====", strconv.Itoa(i), " org===", org)
		for j, orgTemp := range channel.Organizations {
			logger.Info("channel.Organizations j====", strconv.Itoa(j), " orgTemp===", orgTemp)
			if j != i {
				for _, peer := range orgTemp.orgPeers {
					channel.Channel.RemovePeer(peer)
					logger.Info("channel.Channel.RemovePeer(peer) peer====", peer)
				}
			}
		}

		peerConfig, err := client.Client.Config().PeersConfig(org.OrgID)
		Url := ""
		for _, p := range peerConfig {
			if p.Primary {
				Url = p.Host + ":" + strconv.Itoa(p.Port)
			}
		}
		for _, peer := range org.orgPeers {
			channel.Channel.AddPeer(peer)
			if peer.URL() == Url {
				channel.Channel.SetPrimaryPeer(peer)
			}
		}

		if err = admin.JoinChannel(client.Client, org.AdminUser, channel.Channel); err != nil {
			return fmt.Errorf("JoinChannel returned error: %v", err)
		}
	}

	//}
	return nil
}

func getPrimaryPeerOrg(channel *Channel) *Organization {
	index := -1
	for i, org := range channel.Organizations {
		if org.OrgID == channel.PrimaryOrgID {
			index = i
		}
	}
	if index == -1 {
		logger.Info("return Organization{}")
		return &Organization{}
	}
	return channel.Organizations[index]
}

// InitConfig ...
func (appClient *APPClient) InitConfig() (apiconfig.Config, error) {
	configImpl, err := config.InitConfig(appClient.ConfigFile)
	if err != nil {
		return nil, err
	}
	return configImpl, nil
}

// InstantiateCC ...
func (channel *Channel) InstallAndInstantiateCC(args []string) error {
	client := channel.AppClient
	primaryOrg := getPrimaryPeerOrg(channel)
	client.Client.SetUserContext(primaryOrg.AdminUser)
	defer client.Client.SetUserContext(primaryOrg.NormalUser)
	chaincodeQueryResponse, err := client.Client.QueryInstalledChaincodes(channel.Channel.PrimaryPeer())
	if err != nil {
		logger.Errorf("QueryInstalledChaincodes return error: %v", err)
	}
	ccFound := false
	if chaincodeQueryResponse != nil {
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			if chaincode.Name == channel.ChainCode.ChainCodeId&& chaincode.Path == channel.ChainCode.ChainCodePath && chaincode.Version == channel.ChainCode.ChainCodeVersion {
				fmt.Printf("Found chaincode: %s\n", chaincode)
				ccFound = true
			}
		}
	}

	if !ccFound {
		// ###### install #######
		if err := channel.InstallCC(nil); err != nil {
			logger.Errorf("installCC return error: %v", err)
		}

		// ###### instantiate #######
		if err := channel.InstantiateCC(nil); err != nil {
			logger.Errorf("instantiateCC return error: %v", err)
		}
	}
	return nil
}

// InstantiateCC ...
func (channel *Channel) InstantiateCC(args []string) error {
	client := channel.AppClient
	primaryOrg := getPrimaryPeerOrg(channel)
	// InstantiateCC requires AdminUser privileges so setting user context with Admin User
	client.Client.SetUserContext(primaryOrg.AdminUser)
	// must reset client user context to normal user once done with Admin privilieges
	defer client.Client.SetUserContext(primaryOrg.NormalUser)

	if err := admin.SendInstantiateCC(channel.Channel, channel.ChainCode.ChainCodeId, args, channel.ChainCode.ChainCodePath, channel.ChainCode.ChainCodeVersion, []apitxn.ProposalProcessor{channel.Channel.PrimaryPeer()}, client.EventHub); err != nil {
		return err
	}
	return nil
}

// InstallCC ...
func (channel *Channel) InstallCC(chaincodePackage []byte) error {
	client := channel.AppClient
	primaryOrg := getPrimaryPeerOrg(channel)
	// installCC requires AdminUser privileges so setting user context with Admin User
	client.Client.SetUserContext(primaryOrg.AdminUser)
	// must reset client user context to normal user once done with Admin privilieges
	defer client.Client.SetUserContext(primaryOrg.NormalUser)

	if err := admin.SendInstallCC(client.Client, channel.ChainCode.ChainCodeId, channel.ChainCode.ChainCodePath, channel.ChainCode.ChainCodeVersion, chaincodePackage, channel.Channel.Peers(), GetDeployPath()); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}

	return nil
}

// GetDeployPath ..
func GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../chaincode")
}

// Invoke ...
func (context *CurrentContext) Invoke(args []string) (string, error) {
	client := context.AppClient
	context.AppClient.Client.SetUserContext(context.Organization.NormalUser)
	fcn := "invoke"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")
	transactionProposalResponse, txID, err := context.CreateAndSendTransactionProposal(fcn, args, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}

	// Register for commit event
	done, fail := client.RegisterTxEvent(txID)

	txResponse, err := context.CreateAndSendTransaction(transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	fmt.Println(txResponse)
	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 100):

		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}
	return txID.ID, nil
}

// RegisterTxEvent registers on the given eventhub for the give transaction
// returns a boolean channel which receives true when the event is complete
// and an error channel for errors
// TODO - Duplicate
func (appClient *APPClient) RegisterTxEvent(txID apitxn.TransactionID) (chan bool, chan error) {
	done := make(chan bool)
	fail := make(chan error)

	appClient.EventHub.RegisterTxEvent(txID, func(txId string, errorCode peer.TxValidationCode, err error) {
		if err != nil {
			fmt.Printf("Received error event for txid(%s)\n", txId)
			fail <- err
		} else {
			fmt.Printf("Received success event for txid(%s)\n", txId)
			done <- true
		}
	})

	return done, fail
}

// CreateAndSendTransactionProposal ... TODO duplicate
func (context *CurrentContext) CreateAndSendTransactionProposal(
fcn string, args []string, transientData map[string][]byte) ([]*apitxn.TransactionProposalResponse, apitxn.TransactionID, error) {
	request := apitxn.ChaincodeInvokeRequest{
		Targets:      []apitxn.ProposalProcessor{context.Channel.Channel.PrimaryPeer()},
		Fcn:          fcn,
		Args:         args,
		TransientMap: transientData,
		ChaincodeID:  context.Channel.ChainCode.ChainCodeId,
	}
	transactionProposalResponses, txnID, err := context.Channel.Channel.SendTransactionProposal(request)
	if err != nil {
		return nil, txnID, err
	}

	for _, v := range transactionProposalResponses {
		if v.Err != nil {
			return nil, txnID, fmt.Errorf("invoke Endorser %s returned error: %v", v.Endorser, v.Err)
		}
	}

	return transactionProposalResponses, txnID, nil
}

// CreateAndSendTransaction ...
func (context *CurrentContext) CreateAndSendTransaction(resps []*apitxn.TransactionProposalResponse) (*apitxn.TransactionResponse, error) {
	tx, err := context.Channel.Channel.CreateTransaction(resps)
	if err != nil {
		return nil, fmt.Errorf("CreateTransaction return error: %v", err)
	}

	transactionResponse, err := context.Channel.Channel.SendTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("SendTransaction return error: %v", err)

	}

	if transactionResponse.Err != nil {
		return nil, fmt.Errorf("Orderer %s return error: %v", transactionResponse.Orderer, transactionResponse.Err)
	}

	return transactionResponse, nil
}

// Query ...
func (context *CurrentContext) Query(args []string) (string, error) {
	context.AppClient.Client.SetUserContext(context.Organization.NormalUser)
	client := context.AppClient
	fcn := "invoke"
	return fabrictxn.QueryChaincode(client.Client, context.Channel.Channel, context.Channel.ChainCode.ChainCodeId, fcn, args)
}
