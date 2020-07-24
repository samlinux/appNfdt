/**
 *
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type DataItem struct {
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Owner       Owner  `json:"owner,omitempty"`
}

type Owner struct {
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	Departement string `json:"departement,omitempty"`
}

// ===========================
// Main
// ===========================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ===========================
// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ===========================
// Invoke central ctrl of the chaincode
// ===========================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// this is printed to the chaincode container
	fmt.Println("nfdt01 Invoke")

	// import function name and arguments
	function, args := stub.GetFunctionAndParameters()

	// ctrol the flow
	if function == "add" {
		return t.add(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "queryById" {
		return t.queryById(stub, args)
	} else if function == "queryByOwner" {
		return t.queryByOwner(stub, args)
	} else if function == "queryAdHoc" {
		return t.queryAdHoc(stub, args)
	} else if function == "getAllTxByKey" {
		return t.getAllTxByKey(stub, args)
	}

	// if no case match an error will be thrown
	return shim.Error("Invalid invoke function name. Expecting \"add\" \"queryById\" \"queryByOwner\"  \"queryAdHoc\"  \"getAllTxByKey\" ")
}

// =====================================
// add a new DataItem to the blockchain
// =====================================
func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// representing an error
	var err error

	// we need two params a key(Asset) and a value
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// get TxId
	txId := stub.GetTxID()

	// the asset key
	key := args[0]

	fmt.Printf("%+s\n", key)

	// build our data object
	jsonObject := args[1]

	// convert the json object to a byte array
	data := []byte(jsonObject)

	fmt.Printf("%+s\n", jsonObject)

	// create the new DataItem object
	var item DataItem
	item = DataItem{}

	// convert the byte array to a struct object
	err1 := json.Unmarshal(data, &item)
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	// debug the hole object
	fmt.Printf("%+v\n", item)

	// convert the struct object return to bytes
	dataItemJSONasBytes, err := json.Marshal(item)
	if err != nil {
		return shim.Error(err.Error())
	}

	// write the state back to the ledger
	err = stub.PutState(key, dataItemJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// we return the new created asset key to the caller
	return shim.Success(t.getKeyAsBytes(key, txId))
}

// =============================================
// update an existing DataItem in the blockchain
// =============================================
func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// representing an error
	var err error

	// we need two params a key(Asset) and a value
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// get TxId
	txId := stub.GetTxID()

	// build our data object
	dataItemId := args[0]
	jsonObject := args[1]

	// get the corresponding datarecord
	dataItemAsBytes, _ := stub.GetState(dataItemId)

	// convert the record into the data struct
	item := DataItem{}
	json.Unmarshal(dataItemAsBytes, &item)

	// parse the input data (json) into a byte array
	data := []byte(jsonObject)

	// update the data struct with the input data
	json.Unmarshal(data, &item)

	// debug
	fmt.Printf("%+v\n", item)

	// parse the data struct into a byte array
	dataItemJSONasBytes, err := json.Marshal(item)
	if err != nil {
		return shim.Error(err.Error())
	}

	// write the state back to the ledger
	err = stub.PutState(dataItemId, dataItemJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// we return the asset Id and the new TxId to the caller
	return shim.Success(t.getKeyAsBytes(dataItemId, txId))
}

// ==============================================
// queryById a DataItem from the blockchain
// ==============================================
func (t *SimpleChaincode) queryById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	// import the uuid and create the querystring
	uuid := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\"}}", uuid)

	// do the query
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ==============================================
// queryByOwner a DataItem from the blockchain
// ==============================================
func (t *SimpleChaincode) queryByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	// import the uuid and create the querystring
	owner := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"owner.lastname\":\"%s\"}}", owner)

	// do the query
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===============================================
// Ad hoc rich query with idividual query string
// ===============================================
func (t *SimpleChaincode) queryAdHoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// import the uuid and create the querystring
	queryString := args[0]

	// do the query
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getAllTxByKeyx^ returns the history to a given key
// =========================================================================================
func (t *SimpleChaincode) getAllTxByKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	fmt.Printf("- start getKeyHistory: %s\n", key)

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
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

		//buffer.WriteString("\"")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		//buffer.WriteString("\"")

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

	//fmt.Printf("- getKeyHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// =========================================================================================
// getQueryResultForQueryStringWithPagination executes the passed in query string with
// pagination info. Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================
// getKeyAsBytes, is a helper function to give the asset key
// back after the inital storage of the asset
// ========================================================
func (t *SimpleChaincode) getKeyAsBytes(key string, txId string) []byte {
	// we construct a new buffer for the output
	// we want finally a json object like {'Key':'uuid', 'TxId':'txId'}

	var buffer bytes.Buffer
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"")
	buffer.WriteString(key)
	buffer.WriteString("\",")
	buffer.WriteString("\"TxId\":")
	buffer.WriteString("\"")
	buffer.WriteString(txId)
	buffer.WriteString("\"")
	buffer.WriteString("}")

	return buffer.Bytes()
}
