package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PolicyReq struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyInfo PolicyInfo `protobuf:"bytes,2,opt,name=policyInfo" json:"policyInfo"`
	OrderInfo  OrderInfo `protobuf:"bytes,3,opt,name=orderInfo" json:"orderInfo"`
}

type PolicyInfo struct {
	PolicyNo           string `protobuf:"bytes,1,opt,name=policyNo" json:"policyNo"`
	InsuranceBeginDate string `protobuf:"bytes,2,opt,name=insuranceBeginDate" json:"insuranceBeginDate"`
	InsuranceEndDate   string `protobuf:"bytes,3,opt,name=insuranceEndDate" json:"insuranceEndDate"`
}

type OrderInfo struct {
	OrderNo               string `protobuf:"bytes,1,opt,name=orderNo" json:"orderNo"`
	ProductCode           string `protobuf:"bytes,2,opt,name=productCode" json:"productCode"`
	PlanCode              string `protobuf:"bytes,3,opt,name=planCode" json:"planCode"`
	PackageCode           string `protobuf:"bytes,4,opt,name=packageCode" json:"packageCode"`
	InsurantName          string `protobuf:"bytes,5,opt,name=insurantName" json:"insurantName"`
	InsurantPhone         string `protobuf:"bytes,6,opt,name=insurantPhone" json:"insurantPhone"`
	InsurantCertificateNo string `protobuf:"bytes,7,opt,name=insurantCertificateNo" json:"insurantCertificateNo"`
}

type FlightReq struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyNo   string `protobuf:"bytes,2,opt,name=policyNo" json:"policyNo"`
	FlightInfo []FlightInfo `protobuf:"bytes,3,opt,name=flightInfo" json:"flightInfo"`
}

type FlightInfo struct {
	FlightNo            string `protobuf:"bytes,1,opt,name=flightNo" json:"flightNo"`
	SectorNo            string `protobuf:"bytes,2,opt,name=sectorNo" json:"sectorNo"`
	PlanedDepartureTime string `protobuf:"bytes,3,opt,name=planedDepartureTime" json:"planedDepartureTime"`
	PlanedArrivalTime   string `protobuf:"bytes,4,opt,name=planedArrivalTime" json:"planedArrivalTime"`
}

type ClaimReq struct {
	OpenId    string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyNo  string `protobuf:"bytes,2,opt,name=policyNo" json:"policyNo"`
	ClaimInfo []ClaimInfo `protobuf:"bytes,3,opt,name=claimInfo" json:"claimInfo"`
}

type ClaimInfo struct {
	FlightNo            string `protobuf:"bytes,1,opt,name=flightNo" json:"flightNo"`
	SectorNo            string `protobuf:"bytes,2,opt,name=sectorNo" json:"sectorNo"`
	PlanedDepartureTime string `protobuf:"bytes,3,opt,name=planedDepartureTime" json:"planedDepartureTime"`
	PlanedArrivalTime   string `protobuf:"bytes,4,opt,name=planedArrivalTime" json:"planedArrivalTime"`
	ActualDepartureTime string `protobuf:"bytes,5,opt,name=actualDepartureTime" json:"actualDepartureTime"`
	ActualArrivalTime   string `protobuf:"bytes,6,opt,name=actualArrivalTime" json:"actualArrivalTime"`
	FlightState         int `protobuf:"bytes,7,opt,name=flightState" json:"flightState"`
	ClaimAmount         float64 `protobuf:"bytes,8,opt,name=claimAmount" json:"claimAmount"`
	ClaimPaymentState   int `protobuf:"bytes,9,opt,name=claimPaymentState" json:"claimPaymentState"`
}

// InsuranceChaincode mutual insurance Chaincode implementation
type InsuranceChaincode struct {
}

// Init ...
func (t *InsuranceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### insurance Init ###########")
	_, _ = stub.GetFunctionAndParameters()

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			fmt.Printf("Transient data in 'init' : %s\n", transientData)
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Invoke ...
func (t *InsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### insurance Invoke ###########")
	function, args := stub.GetFunctionAndParameters()
	switch args[0] {
	case "postPolicy":
		return t.postPolicy(stub, args)
	case "getPolicy":
		return t.getPolicy(stub, args)
	case "postFlight":
		return t.postFlight(stub, args)
	case "getFlight":
		return t.getFlight(stub, args)
	case "postClaim":
		return t.postClaim(stub, args)
	case "getClaim" :
		return t.getClaim(stub, args)
	default:
		return shim.Error("Unknown function: " + function)
	}
}

// postPolicy postPolicy
func (t *InsuranceChaincode) postPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("postPolicy=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[1])
	var policyReq PolicyReq
	err := json.Unmarshal(body, &policyReq)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("postPolicy: '%v'\n", policyReq)

	attributes := []string{policyReq.PolicyInfo.PolicyNo}
	key, err := stub.CreateCompositeKey("policy", attributes)
	if err != nil {
		return shim.Error(err.Error())
	}

	policyReqByte, err := json.Marshal(policyReq)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, policyReqByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *InsuranceChaincode) getPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("getPolicy=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	key, err := stub.CreateCompositeKey("policy", []string{args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}

	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func (t *InsuranceChaincode) postFlight(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("postFlight=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	body := []byte(args[1])
	var flightReq FlightReq
	err := json.Unmarshal(body, &flightReq)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("postFlight: '%v'\n", flightReq)

	attributes := []string{flightReq.PolicyNo}
	key, err := stub.CreateCompositeKey("flight", attributes)
	if err != nil {
		return shim.Error(err.Error())
	}

	flightReqByte, err := json.Marshal(flightReq)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, flightReqByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *InsuranceChaincode) getFlight(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("getFlight=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}
	key, err := stub.CreateCompositeKey("flight", []string{args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}

	value, err := stub.GetState(string(key))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func (t *InsuranceChaincode) postClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("postClaim=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	body := []byte(args[1])
	var claimReq ClaimReq
	err := json.Unmarshal(body, &claimReq)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("postClaim: '%v'\n", claimReq)

	attributes := []string{claimReq.PolicyNo}
	key, err := stub.CreateCompositeKey("claim", attributes)
	if err != nil {
		return shim.Error(err.Error())
	}

	claimReqByte, err := json.Marshal(claimReq)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, claimReqByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *InsuranceChaincode) getClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("getClaim=%v\n", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}
	key, err := stub.CreateCompositeKey("claim", []string{args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}

	value, err := stub.GetState(string(key))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func main() {
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}