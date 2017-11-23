/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	//"fmt"
	//"math"
	//"strconv"
	//"testing"
	//"time"

	//"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	fabrictxn "github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn"
	"fmt"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	//"math"
	//"time"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	//"github.com/docker/engine-api/types/filters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/hyperledger/fabric/common/flogging"

	"runtime"
	"os"
	logging "github.com/op/go-logging"
	"sync"
	"time"
	"strconv"
	"io/ioutil"
)

const (
	pollRetries = 5
)

var logger = logging.MustGetLogger("performancetest")

var mainCmd = &cobra.Command{
	Use: "app",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loggingSpec := viper.GetString("logging_level")
		if loggingSpec == "" {
			loggingSpec = "DEBUG"
		}
		flogging.InitFromSpec(loggingSpec)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// TestOrgsEndToEnd creates a channel with two organisations, installs chaincode
// on each of them, and finally invokes a transaction on an org2 peer and queries
// the result from an org1 peer
func main() {
	//mainFlags := mainCmd.PersistentFlags()
	//mainFlags.BoolVarP(&versionFlag, "version", "v", false, "Display current version of fabric peer server")
	mainCmd.AddCommand(initCmd())
	mainCmd.AddCommand(invokeCmd())
	mainCmd.AddCommand(queryCmd())
	runtime.GOMAXPROCS(viper.GetInt("core.gomaxprocs"))

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
	logger.Info("Exiting.....")
}

func initCmd() *cobra.Command {
	return initStartCmd
}

func invokeCmd() *cobra.Command {
	return invokeStartCmd
}

func queryCmd() *cobra.Command {
	return queryStartCmd
}

var (
	initStartCmd = &cobra.Command{
		Use:   "init",
		Short: "Starts to test fabric.",
		Long:  `Starts to test the network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			inittest(args)
			return nil
		},
	}

	invokeStartCmd = &cobra.Command{
		Use:   "invoke",
		Short: "Starts to test fabric.",
		Long:  `Starts to test the network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			startoinvoke(args)
			return nil
		},
	}

	queryStartCmd = &cobra.Command{
		Use:   "query",
		Short: "Starts to test fabric.",
		Long:  `Starts to test the network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			startoquery(args)
			return nil
		},
	}
)

func inittest(args []string) {
	// Bootstrap network
	initializeFabricClient()
	loadOrgUsers()
	loadOrgPeers()
	loadOrderer()
	loadTestChannel()
	//joinTestChannel()
	installAndInstantiate()

	fmt.Printf("org1 peer0 is %+v\n", org1TestPeer0)
	fmt.Printf("org2 peer0 is %+v\n", org2TestPeer0)
	fmt.Printf("org3 peer0 is %+v\n", org3TestPeer0)
	fmt.Printf("org4 peer0 is %+v\n", org4TestPeer0)

	//invokeANDquery()
}

func startoinvoke(args []string) {
	// Bootstrap network
	initializeFabricClient()
	loadOrgUsers()
	loadOrgPeers()
	loadOrderer()
	loadTestChannel()

	fmt.Printf("org1 peer0 is %+v\n", org1TestPeer0)
	fmt.Printf("org2 peer0 is %+v\n", org2TestPeer0)
	fmt.Printf("org3 peer0 is %+v\n", org3TestPeer0)
	fmt.Printf("org4 peer0 is %+v\n", org4TestPeer0)

	testmaininvoke(args)
}

func startoquery(args []string) {
	// Bootstrap network
	initializeFabricClient()
	loadOrgUsers()
	loadOrgPeers()
	loadOrderer()
	loadTestChannel()

	fmt.Printf("org1 peer0 is %+v\n", org1TestPeer0)
	fmt.Printf("org2 peer0 is %+v\n", org2TestPeer0)
	fmt.Printf("org3 peer0 is %+v\n", org3TestPeer0)
	fmt.Printf("org4 peer0 is %+v\n", org4TestPeer0)

	testmainquery(args)
}



func testmaininvoke(args []string) error {
	loops := 10
	var alivetime float64 = 1
	if len(args) == 1 {
		if args[0] != "" {
			loop, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				loops = loop
			}
		}
	} else if len(args) == 2 {
		if args[0] != "" {
			loop, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				loops = loop
			}
		}
		if args[1] != "" {
			keepalive, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				alivetime = float64(keepalive)
			}
		}

	}

	var err error
	floc := "./../fabric.perf/"

	fname := "invoke-kplve-" + strconv.FormatFloat(alivetime, 'E', -1, 64) + "-ccrc-" + strconv.Itoa(loops) + "-"+ time.Now().String()+".log"
	pfname := floc + fname
	fabrictxn.Txlogfd, err =os.OpenFile(pfname, os.O_RDWR|os.O_CREATE, 0766);


	if err != nil {
		fmt.Println("Open file failed:", err)
		return err
	}
	defer fabrictxn.Txlogfd.Close()


	resfname := floc + "RES-" + fname
	reslogfd, err := os.OpenFile(resfname, os.O_RDWR|os.O_CREATE, 0766);

	if err != nil {
		fmt.Println("res Open file failed:", err)
		return err
	}
	defer reslogfd.Close()

	alife := time.Now()

	for {
		logger.Infof("Starting test...")
		invokeSuccessNum = 0
		//fmt.Println(loops)

		var wg sync.WaitGroup
		wg.Add(loops)
		t1 := time.Now()
		for i := 1; i <= loops; i ++ {
			if i > 100 && i < 1000 {
				if i % 10 == 0{
				fmt.Print(i)
				fmt.Print(" ")
				}
			} else if i > 1000 {
				if i % 100 == 0{
					fmt.Print(i)
					fmt.Print(" ")
				}
			} else {
				fmt.Print(i)
				fmt.Print(" ")
			}


			go func() {
				defer wg.Done()
				funcname := "initMarble"
				testinvoke(funcname)
			}()
			time.Sleep(time.Duration(int64(time.Second)/int64(loops)))
		}
		fmt.Println()
		wg.Wait()
		elapsed := time.Since(t1)
		life := time.Since(alife)

		invokeResult := fmt.Sprintf("===============test invoke result===============\n" +
			"| keeped time \t| cost time \t| success num \t|\n" +
			"------------------------------------------------\n" +
			"| %s \t| %s \t| %d \t|\n" +
			"------------------------------------------------\n",
			life, elapsed, invokeSuccessNum)
		fmt.Print(invokeResult)
		reslogfd.WriteString(invokeResult)
		break
		if life.Minutes() > alivetime {
			break
		}

	}

	unloadOrgPeersEventHub()

	return nil

}

func generateInvokeArgs() []string {

	index, err := readIndex(fname)
	if err != nil {
		fmt.Println("getIndex err: ", err)
		index = 0
	}
	marblea := fmt.Sprintf("marblea%d", index)
	bluea := fmt.Sprintf("bluea%d", index)
	size := fmt.Sprintf("%d", index + 30)
	toma := fmt.Sprintf("toma%d",index)

	args := []string{marblea, bluea, size, toma }

	return args
}

func testinvoke(funcname string) {

	fcn := funcname

	orgTestClient.SetUserContext(org1User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	_, err := fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org1peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failInvokeTestIfError(err, "InvokeChaincode")

	//fmt.Println("=c1=1=========result after change:", result)
}

func testinvoke2(funcname string) {

	fcn := funcname

	orgTestClient.SetUserContext(org2User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	_, err := fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org2peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failInvokeTestIfError(err, "InvokeChaincode")

	//fmt.Println("=c1=2=========result after change:", result)
}

func testinvoke3(funcname string) {

	fcn := funcname

	orgTestClient.SetUserContext(org3User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	_, err := fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org3peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failInvokeTestIfError(err, "InvokeChaincode")

	//fmt.Println("=c1=3=========result after change:", result)
}

func testinvoke4(funcname string) {

	fcn := funcname

	orgTestClient.SetUserContext(org4User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	_, err := fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org4peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failInvokeTestIfError(err, "InvokeChaincode")

	//fmt.Println("=c1=4=========result after change:", result)
}


func testmainquery(args []string) error {

	loops := 10
	var alivetime float64 = 1
	if len(args) == 1 {
		if args[0] != "" {
			loop, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				loops = loop
			}
		}
	} else if len(args) == 2 {
		if args[0] != "" {
			loop, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				loops = loop
			}
		}
		if args[1] != "" {
			keepalive, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("loop err in args conv: ", err)
			} else {
				alivetime = float64(keepalive)
			}
		}

	}

	var err error
	floc := "./../fabric.perf/"
	fname := floc + "query-kplve-" + strconv.FormatFloat(alivetime, 'E', -1, 64) + "-ccrc-" + strconv.Itoa(loops) + "-"+ time.Now().String()+".log"
	fabrictxn.Txlogfd, err =os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0766);

	if err != nil {
		fmt.Println("Open file failed:", err)
		return err
	}
	defer fabrictxn.Txlogfd.Close()

	alife := time.Now()

	for {
		logger.Infof("Starting query...")
		//fmt.Println(loops)

		var wg sync.WaitGroup
		wg.Add(loops)
		t1 := time.Now()
		for i := 1; i <= loops; i ++ {
			if i > 100 && i < 1000 {
				if i % 10 == 0{
					fmt.Print(i)
					fmt.Print(" ")
				}
			} else if i > 1000 {
				if i % 100 == 0{
					fmt.Print(i)
					fmt.Print(" ")
				}
			} else {
				fmt.Print(i)
				fmt.Print(" ")
			}


			go func() {
				defer wg.Done()
				funcname := "getMarblesByRange"
				testquery(funcname)
				testquery2(funcname)
				testquery3(funcname)
				testquery4(funcname)

			}()
			time.Sleep(time.Duration(int64(time.Second)/int64(loops)))
		}
		fmt.Println()
		wg.Wait()
		elapsed := time.Since(t1)
		fmt.Println("test query elapsed: ", elapsed)

		life := time.Since(alife)
		fmt.Println("run time keeps: ", life)
		break
		if life.Minutes() > alivetime {
			break
		}

	}

	unloadOrgPeersEventHub()

	return nil

}

func generateQueryArgs() []string {
	args := []string{"marblea300", "marblea400"}
	return args
}

func testquery(funcname string) {

	fcn := funcname

	// Query value on org1 peer
	orgTestClient.SetUserContext(org1User)
	orgTestChannel1.SetPrimaryPeer(org1TestPeer0)

	result2, err := fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=p0.g1=========result2 after query:", result2)

}

func testquery2(funcname string) {

	fcn := funcname

	// Query value on org2 peer
	orgTestClient.SetUserContext(org2User)
	orgTestChannel1.SetPrimaryPeer(org2TestPeer0)

	result2, err := fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=p0.g2=========result2 after query:", result2)

}

func testquery3(funcname string) {

	fcn := funcname

	// Query value on org3 peer
	orgTestClient.SetUserContext(org3User)
	orgTestChannel1.SetPrimaryPeer(org3TestPeer0)

	result2, err := fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=p0.g3=========result2 after query:", result2)
}

func testquery4(funcname string) {

	fcn := "getMarblesByRange"

	// Query value on org4 peer
	orgTestClient.SetUserContext(org4User)
	orgTestChannel1.SetPrimaryPeer(org4TestPeer0)

	result2, err := fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=p0.g4=========result2 after query:", result2)
}


func invokeANDquery() {

	fcn := "invoke"

	// Change value on org1 peer with different user
	orgTestClient.SetUserContext(org1User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	result, err := fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org1peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failTestIfError(err, "InvokeChaincode")

	fmt.Println("=c1=1=========result after change:", result)

	orgTestClient.SetUserContext(org2User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	result, err = fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org2peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failTestIfError(err, "InvokeChaincode")

	fmt.Println("=c1=2=========result after change:", result)

	orgTestClient.SetUserContext(org3User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	result, err = fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org3peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failTestIfError(err, "InvokeChaincode")

	fmt.Println("=c1=3=========result after change:", result)

	orgTestClient.SetUserContext(org4User)
	orgTestChannel1.SetPrimaryPeer(*endorsePeer)
	result, err = fabrictxn.InvokeChaincode(orgTestClient, orgTestChannel1, []apitxn.ProposalProcessor{*endorsePeer},
		org4peer0EventHub, gchainCodeID, fcn, generateInvokeArgs(), nil)
	failTestIfError(err, "InvokeChaincode")

	fmt.Println("=c1=4=========result after change:", result)
//////////////////////////////////////////////////////////////////////////////////////////////////
	// Query value on org1 peer
	orgTestClient.SetUserContext(org1User)
	orgTestChannel1.SetPrimaryPeer(org1TestPeer0)

	result2, err := fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=5=========result2 after query:", result2)

	// Query value on org2 peer
	orgTestClient.SetUserContext(org2User)
	orgTestChannel1.SetPrimaryPeer(org2TestPeer0)

	result2, err = fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=6=========result2 after query:", result2)

	// Query value on org3 peer
	orgTestClient.SetUserContext(org3User)
	orgTestChannel1.SetPrimaryPeer(org3TestPeer0)

	result2, err = fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=7=========result2 after query:", result2)

	// Query value on org4 peer
	orgTestClient.SetUserContext(org4User)
	orgTestChannel1.SetPrimaryPeer(org4TestPeer0)

	result2, err = fabrictxn.QueryChaincode(orgTestClient, orgTestChannel1,
		gchainCodeID, fcn, generateQueryArgs())
	failTestIfError(err, "QueryChaincode")
	fmt.Println("=c1=8=========result2 after query:", result2)

	//////////////////////////////////////////////////////////////////////////////////////////////////
}

func unloadOrgPeersEventHub() {

	fmt.Println("disconnect eventhubs")

	orgTestClient.SetUserContext(org1User)
	org1peer0EventHub.Disconnect()

	orgTestClient.SetUserContext(org2User)
	org2peer0EventHub.Disconnect()

	orgTestClient.SetUserContext(org3User)
	org3peer0EventHub.Disconnect()

	orgTestClient.SetUserContext(org4User)
	org4peer0EventHub.Disconnect()

}


var fname = ".valindex"

func writeIndex(fname string, data []byte) error {
	err := ioutil.WriteFile(fname, data, os.ModePerm)
	if err != nil {
		fmt.Println("ioutil write file failed:", err)
		return err
	}
	return nil
}

func readIndex(fname string) (int, error){

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("ioutil read file err:", err)

		writeIndex(fname, []byte("0"))

		return 0, err
	}

	vd, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("loop err in args conv: ", err)
		return 0, err
	}
	vd++
	res := strconv.Itoa(vd)

	err = writeIndex(fname, []byte(res))
	if err != nil {
		fmt.Println("writeIndex err:", err)
		return 0, nil
	}

	return vd, nil
}