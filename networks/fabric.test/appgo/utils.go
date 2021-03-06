package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
)

// GetOrdererAdmin returns a pre-enrolled orderer admin user
func GetOrdererAdmin(c fab.FabricClient, orgName string, networkName string) (ca.User, error) {
	keyDir := fmt.Sprintf("ordererOrganizations/%s.com/users/Admin@%s.com/msp/keystore", networkName, networkName)
	certDir := fmt.Sprintf("ordererOrganizations/%s.com/users/Admin@%s.com/msp/signcerts", networkName, networkName)
	return getDefaultImplPreEnrolledUser(c, keyDir, certDir, "ordererAdmin", orgName)
}

// GetAdmin returns a pre-enrolled org admin user
func GetAdmin(c fab.FabricClient, orgPath string, orgName string, networkName string) (ca.User, error) {
	keyDir := fmt.Sprintf("peerOrganizations/%s.%s.com/users/Admin@%s.%s.com/msp/keystore", orgPath, networkName, orgPath, networkName)
	certDir := fmt.Sprintf("peerOrganizations/%s.%s.com/users/Admin@%s.%s.com/msp/signcerts", orgPath, networkName, orgPath, networkName)
	username := fmt.Sprintf("peer%sAdmin", orgPath)
	return getDefaultImplPreEnrolledUser(c, keyDir, certDir, username, orgName)
}

// GetUser returns a pre-enrolled org user
func GetUser(c fab.FabricClient, orgPath string, orgName string, networkName string) (ca.User, error) {
	keyDir := fmt.Sprintf("peerOrganizations/%s.%s.com/users/User1@%s.%s.com/msp/keystore", orgPath, networkName, orgPath, networkName)
	certDir := fmt.Sprintf("peerOrganizations/%s.%s.com/users/User1@%s.%s.com/msp/signcerts", orgPath, networkName, orgPath, networkName)
	username := fmt.Sprintf("peer%sUser1", orgPath)
	return getDefaultImplPreEnrolledUser(c, keyDir, certDir, username, orgName)
}

// GenerateRandomID generates random ID
func GenerateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return randomString(10)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// GetDefaultImplPreEnrolledUser ...
func getDefaultImplPreEnrolledUser(client fab.FabricClient, keyDir string, certDir string, username string, orgName string) (ca.User, error) {
	privateKeyDir := filepath.Join(client.Config().CryptoConfigPath(), keyDir)
	privateKeyPath, err := getFirstPathFromDir(privateKeyDir)
	if err != nil {
		return nil, fmt.Errorf("Error finding the private key path: %v", err)
	}

	enrollmentCertDir := filepath.Join(client.Config().CryptoConfigPath(), certDir)
	enrollmentCertPath, err := getFirstPathFromDir(enrollmentCertDir)
	if err != nil {
		return nil, fmt.Errorf("Error finding the enrollment cert path: %v", err)
	}
	mspID, err := client.Config().MspID(orgName)
	if err != nil {
		return nil, fmt.Errorf("Error reading MSP ID config: %s", err)
	}
	return deffab.NewPreEnrolledUser(client.Config(), privateKeyPath, enrollmentCertPath, username, mspID, client.CryptoSuite())
}

// Gets the first path from the dir directory
func getFirstPathFromDir(dir string) (string, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Could not read directory %s, err %s", err, dir)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fullName := filepath.Join(dir, string(filepath.Separator), f.Name())
		//fmt.Printf("Reading file %s\n", fullName)
		return fullName, nil
	}

	return "", fmt.Errorf("No paths found in directory: %s", dir)
}

// HasPrimaryPeerJoinedChannel checks whether the primary peer of a channel
// has already joined the channel. It returns true if it has, false otherwise,
// or an error
func HasPrimaryPeerJoinedChannel(client fab.FabricClient, orgUser ca.User, channel fab.Channel) (bool, error) {
	logger.Info("HasPrimaryPeerJoinedChannel start")
	foundChannel := false
	primaryPeer := channel.PrimaryPeer()
	currentUser := client.UserContext()
	defer client.SetUserContext(currentUser)

	client.SetUserContext(orgUser)
	response, err := client.QueryChannels(primaryPeer)
	if err != nil {
		return false, fmt.Errorf("Error querying channel for primary peer: %s", err)
	}
	for _, responseChannel := range response.Channels {
		if responseChannel.ChannelId == channel.Name() {
			foundChannel = true
		}
	}
	logger.Info("HasPrimaryPeerJoinedChannel end")
	return foundChannel, nil
}
