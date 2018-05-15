package main

import (
	"fmt"
	"bytes"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleRegistry struct {
}

func (t *SimpleRegistry) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	A := args[0]
	Aval := args[1]

	// Write the state to the ledger
	err = stub.PutState(A, []byte(Aval))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleRegistry) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addnew" {
		return t.addnew(stub,args)
	} else if function == "transfer" {
		return t.transfer(stub,args)
	} else if function == "query" {
		return t.query(stub,args)
	} else if function == "queryAll" {
		return t.queryAll(stub)
	} else if function == "getHistory" {
		return t.getHistory(stub,args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"addnew\", \"transfer\", \"query\", \"getHistory\" or \"queryAll\"")
}

// func to add a new entity to the chain
func (t *SimpleRegistry) addnew(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Expecting 2 arguments")
	}

	locID := args[0]    		// location ID of land [KEY]
	ownerName := args[1] 		// current owner Name [VALUE]

	// Write the state to the ledger
	err = stub.PutState(locID, []byte(ownerName))  //locID in string and ownerName in bytes 
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}


// transfer owner of the property 
func (t *SimpleRegistry) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Expecting 2 arguments")
	}

	locID := args[0]
	newOwner := args[1]

	_, err = stub.GetState(locID)
	if err != nil {
		return shim.Error("Failed to get record")
	}

	err = stub.PutState(locID, []byte(newOwner))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// callback query function
func (t *SimpleRegistry) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Expecting only location ID for query")
	}

	locID := args[0]

	ownerIDinbytes, err := stub.GetState(locID)
	if err != nil {
		return shim.Error("Failed to get record")
	}

	fmt.Printf("Current Owner for location ID"+ locID +"is"+ string(ownerIDinbytes))
	return shim.Success(ownerIDinbytes)
}

func main() {
	err := shim.Start(new(SimpleRegistry))
	if err != nil {
		fmt.Printf("Error starting Registry: %s", err)
	}
}


// func to query All the existing value in chain
func (t *SimpleRegistry) queryAll(stub shim.ChaincodeStubInterface) pb.Response {

	startKey := "000IND"
	endKey := "999IND"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("------->")
		buffer.WriteString("{Key : ")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(" | ")

		buffer.WriteString("Record : ")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//getting ownership History of the given string key  
func (t *SimpleRegistry) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	keyLand := args[0]

	fmt.Printf("- start getHistory for: %s\n", keyLand)

	resultsIterator, err := stub.GetHistoryForKey(keyLand)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("------->")
		buffer.WriteString("{TransactionID : ")
		buffer.WriteString(response.TxId)
		buffer.WriteString(" | ")

		buffer.WriteString("Value : ")
		buffer.WriteString(string(response.Value))
		buffer.WriteString(" | ")


		buffer.WriteString("Timestamp : ")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}