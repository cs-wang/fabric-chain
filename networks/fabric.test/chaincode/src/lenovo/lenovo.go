package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type InsuranceLedger struct {
	PolicyNo		string `protobuf:"bytes,1,opt,name=PolicyNo" json:"PolicyNo"`
	InsurantID		string `protobuf:"bytes,2,opt,name=InsurantID" json:"InsurantID"`
	ServiceAgreementHASH	string `protobuf:"bytes,3,opt,name=serviceAgreementHASH" json:"serviceAgreementHASH"`
	IMEINo			string `protobuf:"bytes,4,opt,name=IMEINo" json:"IMEINo"`
	ActivateStoreID		string `protobuf:"bytes,5,opt,name=ActivateStoreID" json:"ActivateStoreID"`
	VerifyResult 		string `protobuf:"bytes,6,opt,name=VerifyResult" json:"VerifyResult"`
	SignDate		string `protobuf:"bytes,7,opt,name=SignDate" json:"SignDate"`
	EffecteiveDate		string `protobuf:"bytes,8,opt,name=EffectiveDate" json:"EffectiveDate"`
	ExpirationDate		string `protobuf:"bytes,9,opt,name=ExpirationDate" json:"ExpirationDate"`
	ClaimFlag		string `protobuf:"bytes,10,opt,name=ClaimFlag" json:"ClaimFlag"`
	ClaimDate		string `protobuf:"bytes,11,opt,name=ClaimDate" json:"ClaimDate"`
	ClaimAmount		string `protobuf:"bytes,12,opt,name=ClaimAmount" json:"ClaimAmount"`
	ModifyFlag		string `protobuf:"bytes,13,opt,name=ModifyFlag" json:"ModifyFlag"`
	ExtraData		[]ExtraLedger	`protobuf:"bytes,14,opt,name=ExtraData" json:"ExtraData"`
}

type ExtraLedger struct {
	Key	string 	`protobuf:"bytes,1,opt,name=Key" json:"Key"`
	Data	string 	`protobuf:"bytes,2,opt,name=Data" json:"Data"`
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
	switch function {
	case "postInsurance":
		return t.postInsurance(stub, args)
	case "getInsurance":
		return t.getInsurance(stub, args)
	default:
		return shim.Error("Unknown function: " + function +  " args: " + args[0] )
	}
}

// postPolicy postPolicy
func (t *InsuranceChaincode) postInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Printf("postInsurance=%v\n", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[0])
	var insuranceData InsuranceLedger
	err := json.Unmarshal(body, &insuranceData)
	if err != nil {
		return shim.Error("json.Unmarshal : " + err.Error())
	}

	//fmt.Printf("postInsurance: '%v'\n", insuranceData)

	attributes := []string{insuranceData.PolicyNo}
	key, err := stub.CreateCompositeKey("lenovo", attributes)

	//fmt.Println("postInsurance key: '%v'",key)

	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	idataByte, err := json.Marshal(insuranceData)
	if err != nil {
		return shim.Error("json.Marshal : " + err.Error())
	}
	err = stub.PutState(key, idataByte)
	if err != nil {
		return shim.Error("stub.PutState : " + err.Error())
	}

	return shim.Success(nil)
}

func (t *InsuranceChaincode) getInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Printf("getInsurance=%v\n", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	key, err := stub.CreateCompositeKey("lenovo", []string{args[0]})
	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error("stub.GetState : " + err.Error())
	}

	//fmt.Println(value)
	return shim.Success(value)
}

func main() {
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting lenovo chaincode: %s", err)
	}
}
