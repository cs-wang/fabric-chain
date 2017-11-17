package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"net/http"

	"github.com/gocraft/web"
	"github.com/spf13/viper"
	"crypto/x509"
	"crypto/tls"
)

// --------------- Structs ---------------

type PolicyReq struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyInfo PolicyInfo `protobuf:"bytes,2,opt,name=policyInfo" json:"policyInfo"`
	OrderInfo  OrderInfo `protobuf:"bytes,3,opt,name=orderInfo" json:"orderInfo"`
}

type PolicyRes struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyInfo PolicyInfo `protobuf:"bytes,2,opt,name=policyInfo" json:"policyInfo"`
	OrderInfo  OrderInfo `protobuf:"bytes,3,opt,name=orderInfo" json:"orderInfo"`
	ResultCode string `protobuf:"bytes,4,opt,name=resultCode" json:"resultCode"`
	ErrorDesc  string `protobuf:"bytes,5,opt,name=errorDesc" json:"errorDesc"`
}

type Result struct {
	ResultCode string `protobuf:"bytes,1,opt,name=resultCode" json:"resultCode"`
	ErrorDesc  string `protobuf:"bytes,2,opt,name=errorDesc" json:"errorDesc"`
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

type FlightRes struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyNo   string `protobuf:"bytes,2,opt,name=policyNo" json:"policyNo"`
	FlightInfo []FlightInfo `protobuf:"bytes,3,opt,name=flightInfo" json:"flightInfo"`
	ResultCode string `protobuf:"bytes,4,opt,name=resultCode" json:"resultCode"`
	ErrorDesc  string `protobuf:"bytes,5,opt,name=errorDesc" json:"errorDesc"`
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

type ClaimRes struct {
	OpenId     string `protobuf:"bytes,1,opt,name=openId" json:"openId"`
	PolicyNo   string `protobuf:"bytes,2,opt,name=policyNo" json:"policyNo"`
	ClaimInfo  []ClaimInfo `protobuf:"bytes,3,opt,name=claimInfo" json:"claimInfo"`
	ResultCode string `protobuf:"bytes,4,opt,name=resultCode" json:"resultCode"`
	ErrorDesc  string `protobuf:"bytes,5,opt,name=errorDesc" json:"errorDesc"`
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

// --------------- InsuranceAPP ---------------

// InsuranceAPP defines the Insurance REST service object.
type InsuranceAPP struct {
}

func buildInsuranceRouter() *web.Router {
	//parent.apps
	router := web.New(InsuranceAPP{})

	// Add middleware
	router.Middleware((*InsuranceAPP).setResponseType)
	channels := viper.GetStringMap("channels")
	for currentChannelK, _ := range channels {
		currentChannelPeers := viper.GetStringSlice("channels." + currentChannelK + ".peers")
		for _, currentPeerOrg := range currentChannelPeers {
			key := currentChannelK + currentPeerOrg
			// --- app router ----
			app := router.Subrouter(InsuranceAPP{}, "/v1/" + key)
			app.Post("/policy", (*InsuranceAPP).postPolicy)
			app.Post("/flight", (*InsuranceAPP).postFlight)
			app.Post("/claim", (*InsuranceAPP).postClaim)
			app.Get("/policy/:id", (*InsuranceAPP).getPolicy)
			app.Get("/flight/:id", (*InsuranceAPP).getFlight)
			app.Get("/claim/:id", (*InsuranceAPP).getClaim)
		}
	}

	return router
}

// basicAuthenticate basic authentication
func (s *InsuranceAPP) basicAuthenticate(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	const basicScheme string = "Basic "

	// Confirm the request is sending Basic Authentication credentials.
	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		logger.Errorf("authentication error: scheme=%v", auth)
		return
	}

	// Get the plain-text username and password from the request.
	// The first six characters are skipped - e.g. "Basic ".
	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		logger.Errorf("authentication error: auth=%v", str)
		return
	}

	// Split on the first ":" character only, with any subsequent colons assumed to be part
	// of the password. Note that the RFC2617 standard does not place any limitations on
	// allowable characters in the password.
	creds := bytes.SplitN(str, []byte(":"), 2)

	if len(creds) != 2 {
		logger.Errorf("authentication error: creds=%v", creds)
		return
	}

	user := string(creds[0])
	pass := string(creds[1])

	// TODO: check user and pass

	// Set header for later use
	req.Header.Set("user", user)
	req.Header.Set("pass", pass)
	logger.Infof("basic authentication: user=%v, pass=%v", user, pass)

	next(rw, req)
}

// setResponseType is a middleware function that sets the appropriate response
// headers. Currently, it is setting the "Content-Type" to "application/json" as
// well as the necessary headers in order to enable CORS for Swagger usage.
func (s *InsuranceAPP) setResponseType(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	next(rw, req)
}

func (s *InsuranceAPP) postPolicy(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result
	// Decode the incoming JSON payload
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var policyReq PolicyReq
	err = json.Unmarshal(body, &policyReq)
	if err != nil {
		result.ResultCode = "02"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	logger.Infof("postPolicy: '%s'.\n", policyReq)

	args := []string{"postPolicy", string(body)}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	txID, err := adapter.Invoke(args)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	result.ResultCode = "00"
	result.ErrorDesc = fmt.Sprintf("%v", txID)
	encoder.Encode(result)
	logger.Infof("postPolicy: '%s'\n", txID)
	return
}

func (s *InsuranceAPP) postFlight(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result
	// Decode the incoming JSON payload
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var flightReq FlightReq
	err = json.Unmarshal(body, &flightReq)
	if err != nil {
		result.ResultCode = "02"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	logger.Infof("postFlight: '%s'.\n", flightReq)

	args := []string{"postFlight", string(body)}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	txID, err := adapter.Invoke(args)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	result.ResultCode = "00"
	result.ErrorDesc = fmt.Sprintf("%v", txID)
	encoder.Encode(result)
	logger.Infof("postFlight: '%s'\n", txID)
	return
}

func (s *InsuranceAPP) postClaim(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result
	// Decode the incoming JSON payload
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var claimReq ClaimReq
	err = json.Unmarshal(body, &claimReq)
	if err != nil {
		result.ResultCode = "02"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	logger.Infof("postClaim: '%s'.\n", claimReq)

	args := []string{"postClaim", string(body)}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	txID, err := adapter.Invoke(args)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	result.ResultCode = "00"
	result.ErrorDesc = fmt.Sprintf("%v", txID)
	encoder.Encode(result)
	logger.Infof("postClaim: '%s'\n", txID)
	return
}

func (s *InsuranceAPP) getPolicy(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result

	id := req.PathParams["id"]
	args := []string{"getPolicy", id}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	response, err := adapter.Query(args)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	if response == "" {
		result.ResultCode = "02"
		result.ErrorDesc = "have not such data"
		encoder.Encode(result)
		logger.Infof("getPolicy: '%s'\n", "")
		return
	}

	var policyReq PolicyReq
	err = json.Unmarshal([]byte(response), &policyReq)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var policyRes PolicyRes
	policyRes = PolicyRes{
		OpenId:policyReq.OpenId,
		PolicyInfo:policyReq.PolicyInfo,
		OrderInfo:policyReq.OrderInfo,
		ResultCode:"00",
		ErrorDesc:"success",
	}

	encoder.Encode(policyRes)
	logger.Infof("getPolicy: '%s'\n", response)

	return
}

// GetClaim GetClaim
func (s *InsuranceAPP) getFlight(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result

	id := req.PathParams["id"]
	args := []string{"getFlight", id}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	response, err := adapter.Query(args)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	if response == "" {
		result.ResultCode = "02"
		result.ErrorDesc = "have not such data"
		encoder.Encode(result)
		logger.Infof("getFlight: '%s'\n", "")
		return
	}

	var flightReq FlightReq
	err = json.Unmarshal([]byte(response), &flightReq)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var flightRes FlightRes
	flightRes = FlightRes{
		OpenId:flightReq.OpenId,
		PolicyNo:flightReq.PolicyNo,
		FlightInfo:flightReq.FlightInfo,
		ResultCode:"00",
		ErrorDesc:"success",
	}

	encoder.Encode(flightRes)
	logger.Infof("getFlight: '%s'\n", response)

	return
}

// getClaimm getClaimm
func (s *InsuranceAPP) getClaim(rw web.ResponseWriter, req *web.Request) {
	encoder := json.NewEncoder(rw)
	var result Result

	id := req.PathParams["id"]
	args := []string{"getClaim", id}
	path := req.URL.Path
	adapter := ContextMap[path[4:20]]
	response, err := adapter.Query(args)
	if err != nil {
		result.ResultCode = "01"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	if response == "" {
		result.ResultCode = "02"
		result.ErrorDesc = "have not such data"
		encoder.Encode(result)
		logger.Infof("getClaim: '%s'\n", "")
		return
	}

	var claimReq ClaimReq
	err = json.Unmarshal([]byte(response), &claimReq)
	if err != nil {
		result.ResultCode = "03"
		result.ErrorDesc = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(result)
		logger.Errorf("Error: %s", err)
		return
	}

	var claimRes ClaimRes
	claimRes = ClaimRes{
		OpenId:claimReq.OpenId,
		PolicyNo:claimReq.PolicyNo,
		ClaimInfo:claimReq.ClaimInfo,
		ResultCode:"00",
		ErrorDesc:"success",
	}

	encoder.Encode(claimRes)
	logger.Infof("getClaim: '%s'\n", response)

	return
}

// StartInsuranceServer initializes the REST service and adds the required
// middleware and routes.
func startInsuranceServer() {
	// Initialize the REST service object
	tlsEnabled := viper.GetBool("app.tls.enabled")

	address := viper.GetString("app.address")

	logger.Infof("Initializing the REST service on %s, TLS is %s.", address, (map[bool]string{true: "enabled", false: "disabled"})[tlsEnabled])
	router := buildInsuranceRouter()

	startServerOneByOne(tlsEnabled, viper.GetString("app.address"), router)
}

/**
guolidong:~$ openssl genrsa -out ca.key 2048
Generating RSA private key, 2048 bit long modulus
...+++
.......................+++
unable to write 'random state'
e is 65537 (0x10001)
guolidong:~$ openssl req -x509 -new -nodes -key ca.key -subj "/CN=pingan.com" -days 5000 -out ca.crt
guolidong:~$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus
............+++
................+++
unable to write 'random state'
e is 65537 (0x10001)
guolidong:~$ openssl req -new -key server.key -subj "/CN=www.paic.com" -out server.csr
guolidong:~$ openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000
Signature ok
subject=/CN=www.paic.com
Getting CA Private Key
unable to write 'random state'
guolidong:~$
guolidong:~$
guolidong:~$ openssl rsa -in server.key -out server.key.public
writing RSA key
guolidong:~$ openssl genrsa -out client.key 2048
Generating RSA private key, 2048 bit long modulus
..+++
..................................................................................................................................+++
unable to write 'random state'
e is 65537 (0x10001)
guolidong:~$ openssl req -new -key client.key -subj "/CN=www.paic.com" -out client.csr
guolidong:~$ openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 5000
Signature ok
subject=/CN=www.paic.com
Getting CA Private Key
unable to write 'random state'

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile client.ext -out client.crt -days 5000
 */
func startServerOneByOne(tlsEnabled bool, currentAddress string, router *web.Router) {

	pool := x509.NewCertPool()
	caCertPath := viper.GetString("app.tls.ca.file")

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    currentAddress,
		Handler:  router,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}


	// Start server
	if tlsEnabled {
		err := s.ListenAndServeTLS(viper.GetString("app.tls.cert.file"), viper.GetString("app.tls.key.file"))
		if err != nil {
			logger.Errorf("ListenAndServeTLS: %s", err)
		}
	} else {
		err := http.ListenAndServe(currentAddress, router)
		if err != nil {
			logger.Errorf("ListenAndServe: %s", err)
		}
	}
}

// start serve
func serve(args []string) error {
	// Create and register the REST service if configured
	startInsuranceServer()

	logger.Infof("Starting app...")

	return nil
}
