
package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	//"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
)

type ExtraLedger struct {
	Key	string 	`protobuf:"bytes,1,opt,name=Key" json:"Key"`
	Data	string 	`protobuf:"bytes,2,opt,name=Data" json:"Data"`
}

type PhysicalSeal struct {
	Serie 				string	`protobuf:"bytes,1,opt,name=Serie" json:"Serie"`
	Origanization 			string	`protobuf:"bytes,2,opt,name=Origanization" json:"Origanization"`
	SealId 				string	`protobuf:"bytes,3,opt,name=SealId" json:"SealId"`
	SealType 			string	`protobuf:"bytes,4,opt,name=SealType" json:"SealType"`
	Type 				string	`protobuf:"bytes,5,opt,name=Type" json:"Type"`
	SealName 			string	`protobuf:"bytes,6,opt,name=SealName" json:"SealName"`
	SealFlag	 		string	`protobuf:"bytes,7,opt,name=SealFlag" json:"SealFlag"`
	Applicant			string	`protobuf:"bytes,8,opt,name=Applicant" json:"Applicant"`
	EoaAsn 				string	`protobuf:"bytes,9,opt,name=EoaAsn" json:"EoaAsn"`
	PassAsn 			string	`protobuf:"bytes,10,opt,name=PassAsn" json:"PassAsn"`
	AsSubmitTime 			string	`protobuf:"bytes,11,opt,name=AsSubmitTime" json:"AsSubmitTime"`
	TransferSignTime 		string	`protobuf:"bytes,12,opt,name=TransferSignTime" json:"TransferSignTime"`
	TransferType 			string	`protobuf:"bytes,13,opt,name=TransferType" json:"TransferType"`
	// newly increased proper field
	ApplicationDepartment 		string	`protobuf:"bytes,14,opt,name=ApplicationDepartment" json:"ApplicationDepartment"`
	SealCenterEnterTime 		string	`protobuf:"bytes,15,opt,name=SealCenterEnterTime" json:"SealCenterEnterTime"`
	SealCenterEnterPerson 		string	`protobuf:"bytes,16,opt,name=SealCenterEnterPerson" json:"SealCenterEnterPerson"`
	Transferees 			string	`protobuf:"bytes,17,opt,name=Transferees" json:"Transferees"`
	TransferDepartment 		string	`protobuf:"bytes,18,opt,name=TransferDepartment" json:"TransferDepartment"`
	TransferTime 			string	`protobuf:"bytes,19,opt,name=TransferTime" json:"TransferTime"`
	SealCenterReceiver 		string	`protobuf:"bytes,20,opt,name=SealCenterReceiver" json:"SealCenterReceiver"`
	SealCenterReceiveTime 		string	`protobuf:"bytes,21,opt,name=SealCenterReceiveTime" json:"SealCenterReceiveTime"`
	// destroy record proper field
	DestroyPrintSeal 		string	`protobuf:"bytes,22,opt,name=DestroyPrintSeal" json:"DestroyPrintSeal"`
	ParceInformation 		string	`protobuf:"bytes,23,opt,name=ParceInformation" json:"ParceInformation"`
	SealCenterHandler 		string	`protobuf:"bytes,24,opt,name=SealCenterHandler" json:"SealCenterHandler"`
	SealCenterSendBackTime 		string	`protobuf:"bytes,25,opt,name=SealCenterSendBackTime" json:"SealCenterSendBackTime"`
	OrgUploadDestroyReceipt 	string	`protobuf:"bytes,26,opt,name=OrgUploadDestroyReceipt" json:"OrgUploadDestroyReceipt"`
	DestroyReceiptUploadTime 	string	`protobuf:"bytes,27,opt,name=DestroyReceiptUploadTime" json:"DestroyReceiptUploadTime"`
	// loan return proper field
	ApplicationOrg 			string	`protobuf:"bytes,28,opt,name=ApplicationOrg" json:"ApplicationOrg"`
	ApplicationTime 		string	`protobuf:"bytes,29,opt,name=ApplicationTime" json:"ApplicationTime"`
	loanReason 			string	`protobuf:"bytes,30,opt,name=loanReason" json:"loanReason"`
	loanSignNum 			string	`protobuf:"bytes,31,opt,name=loanSignNum" json:"loanSignNum"`
	SealCenterLoanHandler 		string	`protobuf:"bytes,32,opt,name=SealCenterLoanHandler" json:"SealCenterLoanHandler"`
	SealCenterLoanHandleTime 	string	`protobuf:"bytes,33,opt,name=SealCenterLoanHandleTime" json:"SealCenterLoanHandleTime"`
	SealUseLoanTransfer 		string	`protobuf:"bytes,34,opt,name=SealUseLoanTransfer" json:"SealUseLoanTransfer"`
	WhoOutWithSeal 			string	`protobuf:"bytes,35,opt,name=WhoOutWithSeal" json:"WhoOutWithSeal"`
	WhenOutWithSeal 		string	`protobuf:"bytes,36,opt,name=WhenOutWithSeal" json:"WhenOutWithSeal"`
	SealOutSideTime 		string	`protobuf:"bytes,37,opt,name=SealOutSideTime" json:"SealOutSideTime"`
	SealOutSideSite 		string	`protobuf:"bytes,38,opt,name=SealOutSideSite" json:"SealOutSideSite"`
	SealCenterRemandHandler 	string	`protobuf:"bytes,39,opt,name=SealCenterRemandHandler" json:"SealCenterRemandHandler"`
	SealCenterRemandHandleTime 	string	`protobuf:"bytes,40,opt,name=SealCenterRemandHandleTime" json:"SealCenterRemandHandleTime"`
	SealUseRemandTransfer 		string	`protobuf:"bytes,41,opt,name=SealUseRemandTransfer" json:"SealUseRemandTransfer"`
	// custody record proper field
	SealEnterTime 			string	`protobuf:"bytes,42,opt,name=SealEnterTime" json:"SealEnterTime"`
	SealEnterPerson 		string	`protobuf:"bytes,43,opt,name=SealEnterPerson" json:"SealEnterPerson"`
	SealImageHash 			string	`protobuf:"bytes,44,opt,name=SealImageHash" json:"SealImageHash"`
	ImageUploadTime 		string	`protobuf:"bytes,45,opt,name=ImageUploadTime" json:"ImageUploadTime"`
	ImageUploadPerson 		string	`protobuf:"bytes,46,opt,name=ImageUploadPerson" json:"ImageUploadPerson"`
	ImageUpdateTime 		string	`protobuf:"bytes,47,opt,name=ImageUpdateTime" json:"ImageUpdateTime"`
	ImageUpdatePerson 		string	`protobuf:"bytes,48,opt,name=ImageUpdatePerson" json:"ImageUpdatePerson"`
	SealDisableTime 		string	`protobuf:"bytes,49,opt,name=SealDisableTime" json:"SealDisableTime"`
	AssociateAddSignNum 		string	`protobuf:"bytes,50,opt,name=AssociateAddSignNum" json:"AssociateAddSignNum"`
	AssociateDestroySignNum 	string	`protobuf:"bytes,51,opt,name=AssociateDestroySignNum" json:"AssociateDestroySignNum"`

	ExtraData		[]ExtraLedger	`protobuf:"bytes,52,opt,name=ExtraData" json:"ExtraData"`
}

type PrintSeal struct {
	Serie 				string	`protobuf:"bytes,1,opt,name=Serie" json:"Serie"`
	Origanization 			string	`protobuf:"bytes,2,opt,name=Origanization" json:"Origanization"`
	SealId 				string	`protobuf:"bytes,3,opt,name=SealId" json:"SealId"`
	SealType 			string	`protobuf:"bytes,4,opt,name=SealType" json:"SealType"`
	SealName 			string	`protobuf:"bytes,5,opt,name=SealName" json:"SealName"`
	SealCenterEnterTime 		string	`protobuf:"bytes,6,opt,name=SealCenterEnterTime" json:"SealCenterEnterTime"`
	SealCenterEnterPerson 		string	`protobuf:"bytes,7,opt,name=SealCenterEnterPerson" json:"SealCenterEnterPerson"`
	SealFile 			string	`protobuf:"bytes,8,opt,name=SealFile" json:"SealFile"`
	SealFileUploader 		string	`protobuf:"bytes,9,opt,name=SealFileUploader" json:"SealFileUploader"`
	SealFileUploadTime 		string	`protobuf:"bytes,10,opt,name=SealFileUploadTime" json:"SealFileUploadTime"`
	SealFileInvalidTime	 	string	`protobuf:"bytes,11,opt,name=SealFileInvalidTime" json:"SealFileInvalidTime"`
	SealFileInvalider 		string	`protobuf:"bytes,12,opt,name=SealFileInvalider" json:"SealFileInvalider"`

	ExtraData		[]ExtraLedger	`protobuf:"bytes,13,opt,name=KDData" json:"ExtraData"`
}

type SignatureSeal struct {
	Serie 				string	`protobuf:"bytes,1,opt,name=Serie" json:"Serie"`
	Origanization 			string	`protobuf:"bytes,2,opt,name=Origanization" json:"Origanization"`
	SealId 				string	`protobuf:"bytes,3,opt,name=SealId" json:"SealId"`
	SealType 			string	`protobuf:"bytes,4,opt,name=SealType" json:"SealType"`
	SealName 			string	`protobuf:"bytes,5,opt,name=SealName" json:"SealName"`
	SealCenterEnterTime 		string	`protobuf:"bytes,6,opt,name=SealCenterEnterTime" json:"SealCenterEnterTime"`
	SealCenterEnterPerson 		string	`protobuf:"bytes,7,opt,name=SealCenterEnterPerson" json:"SealCenterEnterPerson"`
	SealFile 			string	`protobuf:"bytes,8,opt,name=SealFile" json:"SealFile"`
	SealFileUploader 		string	`protobuf:"bytes,9,opt,name=SealFileUploader" json:"SealFileUploader"`
	SealFileUploadTime 		string	`protobuf:"bytes,10,opt,name=SealFileUploadTime" json:"SealFileUploadTime"`
	SealFileChanger 		string	`protobuf:"bytes,11,opt,name=SealFileChanger" json:"SealFileChanger"`
	SealFileChangeTime 		string	`protobuf:"bytes,12,opt,name=SealFileChangeTime" json:"SealFileChangeTime"`
	CertFileUploader 		string	`protobuf:"bytes,13,opt,name=CertFileUploader" json:"CertFileUploader"`
	CertFileInvalidTime 		string	`protobuf:"bytes,14,opt,name=CertFileInvalidTime" json:"CertFileInvalidTime"`
	CertFileInvalider 		string	`protobuf:"bytes,15,opt,name=CertFileInvalider" json:"CertFileInvalider"`
	CertType 			string	`protobuf:"bytes,16,opt,name=CertType" json:"CertType"`
	CertFile 			string	`protobuf:"bytes,17,opt,name=CertFile" json:"CertFile"`

	ExtraData		[]ExtraLedger	`protobuf:"bytes,18,opt,name=ExtraData" json:"ExtraData"`
}

type SealRecord struct {
	FileNo 				string	`protobuf:"bytes,1,opt,name=FileNo" json:"FileNo"`
	Title 				string	`protobuf:"bytes,2,opt,name=Title" json:"Title"`
	SealId 				string	`protobuf:"bytes,3,opt,name=SealId" json:"SealId"`
	SealName 			string	`protobuf:"bytes,4,opt,name=SealName" json:"SealName"`
	SealType 			string	`protobuf:"bytes,5,opt,name=SealType" json:"SealType"`
	SourceSystem 			string	`protobuf:"bytes,6,opt,name=SourceSystem" json:"SourceSystem"`
	ApplicationNo 			string	`protobuf:"bytes,7,opt,name=ApplicationNo" json:"ApplicationNo"`
	Applicant 			string	`protobuf:"bytes,8,opt,name=Applicant" json:"Applicant"`
	ApplicationOrg 			string	`protobuf:"bytes,9,opt,name=ApplicationOrg" json:"ApplicationOrg"`
	ApplicationDate 		string	`protobuf:"bytes,10,opt,name=ApplicationDate" json:"ApplicationDate"`
	AttachmentList 			string	`protobuf:"bytes,11,opt,name=AttachmentList" json:"AttachmentList"`
	RequisitionStatus 		string	`protobuf:"bytes,12,opt,name=RequisitionStatus" json:"RequisitionStatus"`
	MainLeaderShip 			string	`protobuf:"bytes,13,opt,name=MainLeaderShip" json:"MainLeaderShip"`
	ApplicationDepartment 		string	`protobuf:"bytes,14,opt,name=ApplicationDepartment" json:"ApplicationDepartment"`
	ManagerPhone 			string	`protobuf:"bytes,15,opt,name=ManagerPhone" json:"ManagerPhone"`
	SealUseSubmitTime 		string	`protobuf:"bytes,16,opt,name=SealUseSubmitTime" json:"SealUseSubmitTime"`
	SealUseFileTile 		string	`protobuf:"bytes,17,opt,name=SealUseFileTile" json:"SealUseFileTile"`
	SealUseFileContent 		string	`protobuf:"bytes,18,opt,name=SealUseFileContent" json:"SealUseFileContent"`
	SealUseType 			string	`protobuf:"bytes,19,opt,name=SealUseType" json:"SealUseType"`
	IsContract 			string	`protobuf:"bytes,20,opt,name=IsContract" json:"IsContract"`
	FaxMethod 			string	`protobuf:"bytes,21,opt,name=FaxMethod" json:"FaxMethod"`
	SealCoordinator 		string	`protobuf:"bytes,22,opt,name=SealCoordinator" json:"SealCoordinator"`
	RequirementsDescribe 		string	`protobuf:"bytes,23,opt,name=RequirementsDescribe" json:"RequirementsDescribe"`
	ReturnAddress 			string	`protobuf:"bytes,24,opt,name=ReturnAddress" json:"ReturnAddress"`
	MailingAddressInformation 	string	`protobuf:"bytes,25,opt,name=MailingAddressInformation" json:"MailingAddressInformation"`
	SealUseRecord 			string	`protobuf:"bytes,26,opt,name=SealUseRecord" json:"SealUseRecord"`
	ApprovalComments 		string	`protobuf:"bytes,27,opt,name=ApprovalComments" json:"ApprovalComments"`
	FileName 			string	`protobuf:"bytes,28,opt,name=FileName" json:"FileName"`
	SealUseTime 			string	`protobuf:"bytes,29,opt,name=SealUseTime" json:"SealUseTime"`
	AttachmentIDBfSealUse 		string	`protobuf:"bytes,30,opt,name=AttachmentIDBfSealUse" json:"AttachmentIDBfSealUse"`
	AttachmentIDAftSealUse 		string	`protobuf:"bytes,31,opt,name=AttachmentIDAftSealUse" json:"AttachmentIDAftSealUse"`
	CertificateType 		string	`protobuf:"bytes,32,opt,name=CertificateType" json:"CertificateType"`
	NumOfSealUse 			string	`protobuf:"bytes,33,opt,name=NumOfSealUse" json:"NumOfSealUse"`

	ExtraData		[]ExtraLedger	`protobuf:"bytes,34,opt,name=ExtraData" json:"ExtraData"`
}

type StamperChaincode struct {

	//ExtraData		[]ExtraLedger	`protobuf:"bytes,1,opt,name=KDData" json:"ExtraData"`
}

func main() {
	err := shim.Start(new(StamperChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (sc *StamperChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (sc *StamperChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "physicalSeal" {
		return sc.physicalSeal(stub, args)
	} else if function == "printSeal" {
		return sc.printSeal(stub, args)
	} else if function == "signatureSeal" {
		return sc.signatureSeal(stub, args)
	} else if function == "querySeal" {
		return sc.querySeal(stub, args)
	} else if function == "sealRecord" {
		return sc.sealRecord(stub, args)
	} else if function == "queryRecord" {
		return sc.queryRecord(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")
}

func (sc *StamperChaincode) physicalSeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke physicalSeal\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[0])
	var physicalSealData PhysicalSeal
	err := json.Unmarshal(body, &physicalSealData)
	if err != nil {
		return shim.Error("json.Unmarshal : " + err.Error())
	}

	//fmt.Printf("postInsurance: '%v'\n", insuranceData)

	attributes := []string{physicalSealData.SealId, physicalSealData.SealType}
	key, err := stub.CreateCompositeKey("seal", attributes)

	//fmt.Println("postInsurance key: '%v'",key)

	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	idataByte, err := json.Marshal(physicalSealData)
	if err != nil {
		return shim.Error("json.Marshal : " + err.Error())
	}
	err = stub.PutState(key, idataByte)
	if err != nil {
		return shim.Error("stub.PutState : " + err.Error())
	}

	return shim.Success(nil)
}

func (sc *StamperChaincode) printSeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke printSeal\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[0])
	var physicalSealData PhysicalSeal
	err := json.Unmarshal(body, &physicalSealData)
	if err != nil {
		return shim.Error("json.Unmarshal : " + err.Error())
	}

	//fmt.Printf("postInsurance: '%v'\n", insuranceData)

	attributes := []string{physicalSealData.SealId, physicalSealData.SealType}
	key, err := stub.CreateCompositeKey("seal", attributes)

	//fmt.Println("postInsurance key: '%v'",key)

	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	idataByte, err := json.Marshal(physicalSealData)
	if err != nil {
		return shim.Error("json.Marshal : " + err.Error())
	}
	err = stub.PutState(key, idataByte)
	if err != nil {
		return shim.Error("stub.PutState : " + err.Error())
	}

	return shim.Success(nil)
}

func (sc *StamperChaincode) signatureSeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke signatureSeal\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[0])
	var signatureSealData SignatureSeal
	err := json.Unmarshal(body, &signatureSealData)
	if err != nil {
		return shim.Error("json.Unmarshal : " + err.Error())
	}

	//fmt.Printf("postInsurance: '%v'\n", insuranceData)

	attributes := []string{signatureSealData.SealId, signatureSealData.SealType}
	key, err := stub.CreateCompositeKey("seal", attributes)

	//fmt.Println("postInsurance key: '%v'",key)

	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	idataByte, err := json.Marshal(signatureSealData)
	if err != nil {
		return shim.Error("json.Marshal : " + err.Error())
	}
	err = stub.PutState(key, idataByte)
	if err != nil {
		return shim.Error("stub.PutState : " + err.Error())
	}

	return shim.Success(nil)
}

func (sc *StamperChaincode) querySeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke querySeal\n")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	key, err := stub.CreateCompositeKey("seal", []string{args[0], args[1]})
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

func (sc *StamperChaincode) sealRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke sealRecord\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// TODO: check duplication
	//txID := stub.GetTxID()

	body := []byte(args[0])
	var sealRecordData SealRecord
	err := json.Unmarshal(body, &sealRecordData)
	if err != nil {
		return shim.Error("json.Unmarshal : " + err.Error())
	}

	//fmt.Printf("postInsurance: '%v'\n", insuranceData)

	attributes := []string{sealRecordData.SealId, sealRecordData.SealType, sealRecordData.FileNo}
	key, err := stub.CreateCompositeKey("seal", attributes)

	//fmt.Println("postInsurance key: '%v'",key)

	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	idataByte, err := json.Marshal(sealRecordData)
	if err != nil {
		return shim.Error("json.Marshal : " + err.Error())
	}
	err = stub.PutState(key, idataByte)
	if err != nil {
		return shim.Error("stub.PutState : " + err.Error())
	}

	return shim.Success(nil)
}

func (sc *StamperChaincode) queryRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("invoke queryRecord\n")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	key, err := stub.CreateCompositeKey("seal", []string{args[0], args[1], args[2]})
	if err != nil {
		return shim.Error("stub.CreateCompositeKey : " + err.Error())
	}

	//value, err := stub.GetState(key)
	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error("stub.GetHistoryForKey : " + err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")


	fmt.Println(buffer.String())
	return shim.Success(buffer.Bytes())
	//return shim.Success(value)
}
