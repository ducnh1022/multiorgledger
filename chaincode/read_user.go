package main

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/multiorgledger/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


func (t *MultiOrgChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println(" ******** Invoke Read User ******** ")

	var user model.User
	var email, role, eventID string
	var queryCreatorOrg string
	var queryCreator string
	var needHistory bool

	email 	= args[1]
	eventID = args[2]
	queryCreatorOrg = args[3]
	needHistory, _ = strconv.ParseBool(args[4])

	role , err := t.getRole(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to get roles from the account: %v", err))
	}

	fmt.Println(" Read User - Role === "+role)

	fmt.Println("##### Read "+email+" User #####")

	indexName := "email"
	userNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{email})

	err = getDataFromLedger(stub, userNameIndexKey, &user)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve userData in the ledger: %v", err))
	}

	userAsByte, err := objectToByte(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the userData to byte: %v", err))
	}

	err = stub.SetEvent(eventID, []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}



		/*	Created History for Read by email Transaction */

		if needHistory {
			if strings.EqualFold(role,model.ADMIN){
				queryCreator = model.GetCustomOrgName(queryCreatorOrg)+" Admin"
			} else {
				queryCreator = email
			}

			query   := args[0]
			remarks := queryCreator+" read "+email+" 's user details"
			t.createHistory(stub, queryCreator, queryCreatorOrg, email, query, remarks)
		}

	
	
	return shim.Success(userAsByte)
}

func (t *MultiOrgChaincode) readAllUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("##### Read All User #####")

	role , err := t.getRole(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to get roles from the account: %v", err))
	}

	fmt.Println(" Read All User - Role : "+role)
	
	if !strings.EqualFold(role,model.ADMIN){
		return shim.Error(fmt.Sprintf("Only admin can read all the user data from the ledger: %v", err))
	}

	var eventID string

	eventID = args[1]

	indexName := "email"

	iterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve the list of resource in the ledger: %v", err))
	}

	allUsers := make([]model.User, 0)

	for iterator.HasNext() {
		keyValueState, errIt := iterator.Next()
		if errIt != nil {
			return shim.Error(fmt.Sprintf("Unable to retrieve a user in the ledger: %v", errIt))
		}
		var user model.User
		err = byteToObject(keyValueState.Value, &user)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to convert a user: %v", err))
		}

		fmt.Println("Read User : "+user.Name+" -- "+user.Email)

		allUsers = append(allUsers, user)
	}

	allUsersAsByte, err := objectToByte(allUsers)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to convert the users list to byte: %v", err))
	}

	err = stub.SetEvent(eventID, []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(allUsersAsByte)
}