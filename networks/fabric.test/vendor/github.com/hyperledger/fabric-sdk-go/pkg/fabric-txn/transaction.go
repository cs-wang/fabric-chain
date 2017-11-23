/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabrictxn

import (
	"fmt"
	"time"
	"strconv"

	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	internal "github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/internal"

	"github.com/op/go-logging"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"os"
	"sync"
)

var logger = logging.MustGetLogger("fabric_sdk_go")


//QueryChaincode ...
func QueryChaincode(client fab.FabricClient, channel fab.Channel, chaincodeID string, fcn string, args []string) (string, error) {
	err := checkCommonArgs(client, channel, chaincodeID)
	if err != nil {
		return "", err
	}

	transactionProposalResponses, _, err := internal.CreateAndSendTransactionProposal(channel,
		chaincodeID, fcn, args, []apitxn.ProposalProcessor{channel.PrimaryPeer()}, nil)

	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal returned error: %v", err)
	}

	return string(transactionProposalResponses[0].ProposalResponse.GetResponse().Payload), nil
}

type TransactionCostLog struct {
	OrderType string  `protobuf:"bytes,1,opt,name=orderType" json:"orderType"`
	Times     float64 `protobuf:"bytes,2,opt,name=times" json:"times"`
	Timestamp string  `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp"`
	TxId      string `protobuf:"bytes,4,opt,name=txId" json:"txId"`
	TypeId    string `protobuf:"bytes,5,opt,name=typeId" json:"typeId"`
}

func postOnceTxLog(txlogs []TransactionCostLog) {

	var resp *http.Response
	var err error
	jsont, err := json.Marshal(txlogs)
	if err != nil {
		fmt.Println("========json.Marshal========")
		fmt.Printf("%+v", txlogs)
		fmt.Println("========err========")
		fmt.Printf("%+v", jsont)
		fmt.Println("================:", err)
		return
	}
	//POST /api/v1/influxdb/insertTranStaticsDat
	//http://10.25.50.50:8099/baas/swagger-ui.html#!/influxdb-controller/insertTranStaticsDataUsingPOST
	//http://10.25.50.50:8099/baas/api/v1/influxdb/insertTranStaticsData
	//http://10.25.50.52:8083/
	//select * from "tran_statics_data"
	//http://10.25.50.51:3000/dashboard/db/influxdata
	resp, err = http.Post("http://10.25.50.50:8099/baas/api/v1/influxdb/insertTranStaticsData", "application/json", bytes.NewReader(jsont))
	//resp, err = httpsClient.Post("https://localhost:8088/v1/channel1peerorg1/insertTranData", "application/json", bytes.NewReader(jsont))

	if err != nil {
		// 处理错误
		fmt.Println("postOnceTxLog Error: err != nil")
		return
	}
	if resp.StatusCode != http.StatusOK {
		// 处理错误
		fmt.Println("postOnceTxLog Error: resp.StatusCode != http.StatusOK: ", resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logger.Info(string(body))
}

var Txlogfd *os.File

var filelock sync.RWMutex

func savefileTxLog(txlogs []TransactionCostLog) {
	filelock.Lock()
	defer filelock.Unlock()
	if Txlogfd == nil {
		fmt.Println("file not opened")
		return
	}
	jsont, err := json.Marshal(txlogs)
	if err != nil {
		fmt.Println("jsont, err := json.Marshal(txlogs), err", err)
		return
	}
	Txlogfd.Write(jsont)
	Txlogfd.Write([]byte("\n"))
}

// InvokeChaincode ...
func InvokeChaincode(client fab.FabricClient, channel fab.Channel, targets []apitxn.ProposalProcessor,
	eventHub fab.EventHub, chaincodeID string, fcn string, args []string, transientData map[string][]byte) (apitxn.TransactionID, error) {

	err := checkCommonArgs(client, channel, chaincodeID)
	if err != nil {
		return apitxn.TransactionID{}, err
	}

	if eventHub == nil {
		return apitxn.TransactionID{}, fmt.Errorf("Eventhub is nil")
	}

	if targets == nil || len(targets) == 0 {
		return apitxn.TransactionID{}, fmt.Errorf("No target peers")
	}

	if eventHub.IsConnected() == false {
		err = eventHub.Connect()
		if err != nil {
			return apitxn.TransactionID{}, fmt.Errorf("Error connecting to eventhub: %v", err)
		}
		defer eventHub.Disconnect()
	}
	/////////////////////////////////////////////////////////////
	blogs := make([]TransactionCostLog, 0)
	t1 := time.Now()

	transactionProposalResponses, txID, err := internal.CreateAndSendTransactionProposal(channel,
		chaincodeID, fcn, args, targets, transientData)

	elapsed1 := time.Since(t1)
	value := TransactionCostLog{
		OrderType: "kafka",
		Times: float64(elapsed1/time.Millisecond)+float64(elapsed1%time.Millisecond)/1e6,
		Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		TxId: txID.ID,
		TypeId: "t1",
	}

	blogs = append(blogs, value)
	//fmt.Println("---------------------------test internal.CreateAndSendTransactionProposal elapsed: ", elapsed1)
	/////////////////////////////////////////////////////////////
	if err != nil {
		return apitxn.TransactionID{}, fmt.Errorf("CreateAndSendTransactionProposal returned error: %v", err)
	}
	/////////////////////////////////////////////////////////////
	t1 = time.Now()

	done, fail := internal.RegisterTxEvent(txID, eventHub)

	elapsed2 := time.Since(t1)
	value2 := TransactionCostLog{
		OrderType: "kafka",
		Times: float64(elapsed2/time.Millisecond)+float64(elapsed2%time.Millisecond)/1e6,
		Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		TxId: txID.ID,
		TypeId: "t2",
	}
	blogs = append(blogs, value2)
	//fmt.Println("---------------------------test internal.RegisterTxEvent elapsed: ", elapsed2)
        /////////////////////////////////////////////////////////////
	t1 = time.Now()

	_, err = internal.CreateAndSendTransaction(channel, transactionProposalResponses)

	elapsed3 := time.Since(t1)
	value3 := TransactionCostLog{
		OrderType: "kafka",
		Times: float64(elapsed3/time.Millisecond)+float64(elapsed3%time.Millisecond)/1e6,
		Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		TxId: txID.ID,
		TypeId: "t3",
	}
	blogs = append(blogs, value3)

	//fmt.Println("---------------------------test internal.CreateAndSendTransaction elapsed: ", elapsed3)
	/////////////////////////////////////////////////////////////
	if err != nil {
		return txID, fmt.Errorf("CreateAndSendTransaction returned error: %v", err)
	}

	select {
	case <-done:
	case err := <-fail:
		return txID, fmt.Errorf("invoke Error received from eventhub for txid(%s), error(%v)", txID, err)
	case <-time.After(time.Second * 1000):
		return txID, fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	//postOnceTxLog(blogs)
	savefileTxLog(blogs)
	return txID, nil
}

// checkCommonArgs ...
func checkCommonArgs(client fab.FabricClient, channel fab.Channel, chaincodeID string) error {
	if client == nil {
		return fmt.Errorf("Client is nil")
	}

	if channel == nil {
		return fmt.Errorf("Channel is nil")
	}

	if chaincodeID == "" {
		return fmt.Errorf("ChaincodeID is empty")
	}

	return nil
}

// RegisterCCEvent registers chain code event on the given eventhub
// @returns {chan bool} channel which receives true when the event is complete
// @returns {object} ChainCodeCBE object handle that should be used to unregister
func RegisterCCEvent(chainCodeID string, eventID string, eventHub fab.EventHub) (chan bool, *fab.ChainCodeCBE) {
	done := make(chan bool)

	// Register callback for CE
	rce := eventHub.RegisterChaincodeEvent(chainCodeID, eventID, func(ce *fab.ChaincodeEvent) {
		logger.Debugf("Received CC event: %v\n", ce)
		done <- true
	})

	return done, rce
}
