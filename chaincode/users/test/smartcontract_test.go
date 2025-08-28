package test

import (
	"users/smartcontract"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

var Stub *shimtest.MockStub
var Scc *contractapi.ContractChaincode

var user1 smartcontract.User = smartcontract.User{
	Name:   "New John Lee",
	Id_key: "A222333444",
	Email:  "AAAA@bbb.com",
}

var user2 smartcontract.User = smartcontract.User{
	Name:   "Benny Cheng",
	Id_key: "B444555666",
	Email:  "CCCCC@bbb.com",
}

var transaction1 smartcontract.Transaction = smartcontract.Transaction{
	Hash:          "0x000000001",
	Amount:        "200",
	Currency_type: "USD",
	Create_at:     "1735349527123456789",
}

var transaction2 smartcontract.Transaction = smartcontract.Transaction{
	Hash:          "0x000000002",
	Amount:        "300",
	Currency_type: "NTD",
	Create_at:     "1735349527123456789",
}

var transaction3 smartcontract.Transaction = smartcontract.Transaction{
	Hash:          "0x000000003",
	Amount:        "600",
	Currency_type: "GBP",
	Create_at:     "1735349527123456789",
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	log.SetOutput(ioutil.Discard)
}

func NewStub() {
	Scc, err := contractapi.NewChaincode(new(smartcontract.SmartContract))
	if err != nil {
		log.Println("NewChaincode failed", err)
		os.Exit(0)
	}
	Stub = shimtest.NewMockStub("main", Scc)
	MockInitLedger()
}

func Test_getUserList(t *testing.T) {
	fmt.Println("MockGetUserList-----------------")
	NewStub()

	MockCreateUser(user1.Id_key, user1.Name, user1.Email)

	MockCreateUser(user2.Id_key, user2.Name, user2.Email)

	result1, err := MockCreateTransaction(user1.Id_key, transaction1.Hash, transaction1.Amount, transaction1.Currency_type, transaction1.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction1", result1)

	users, err := MockGetUserList()
	if err != nil {
		fmt.Println("GetUserList error", err)
	}
	fmt.Println(users)

	assert.Equal(t, len(users), 3)

}

func Test_CreateUser(t *testing.T) {
	fmt.Println("Test_CreateUser-----------------")
	NewStub()

	err := MockCreateUser(user1.Id_key, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	//這邊會因為重複新增同一個人而報錯
	err2 := MockCreateUser(user1.Id_key, user1.Name, user1.Email)
	if err2 != nil {
		t.FailNow()
	}

}

func Test_UpdateUser(t *testing.T) {
	fmt.Println("Test_UpdateUser-----------------")
	NewStub()

	err := MockCreateUser(user1.Id_key, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	BfUsers, err := MockGetUserList()
	if err != nil {
		fmt.Println("GetUserList error", err)
	}
	fmt.Println("BEFORE ::", BfUsers)

	newName := "Ariel AAA"
	newEmail := "IhaveChanged@com"

	MockUpdateUser(user1.Id_key, newName, newEmail)

	AfUsers, err := MockGetUserList()
	fmt.Println("AFTER ::", AfUsers)

	userJson, err := MockGetUserAndTransactions(user1.Id_key)

	assert.Equal(t, userJson.Id_key, user1.Id_key)
	assert.Equal(t, userJson.Name, newName)
	assert.Equal(t, userJson.Email, newEmail)

}

func Test_TransactionHashExist(t *testing.T) {
	fmt.Println("Test_TransactionHashExist-----------------")
	NewStub()

	err := MockCreateUser(user1.Id_key, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	result1, err := MockCreateTransaction(user1.Id_key, transaction1.Hash, transaction1.Amount, transaction1.Currency_type, transaction1.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction1", result1)

	result2, err := MockCreateTransaction(user1.Id_key, transaction2.Hash, transaction2.Amount, transaction2.Currency_type, transaction2.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}

	fmt.Println("CreateTransaction transaction2", result2)

	result3, err := MockCreateTransaction(user1.Id_key, transaction3.Hash, transaction3.Amount, transaction3.Currency_type, transaction3.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction2", result3)

	TransactionCheckResult, err := MockGetUserByTransactionHash(transaction1.Hash)
	fmt.Print(TransactionCheckResult)

	falseHash := "AABBCCD"
	TransactionCheckResult2, err := MockGetUserByTransactionHash(falseHash)
	if err != nil {
		fmt.Println("get TransactionCheckResult2", falseHash, "error", err)
		fmt.Print(TransactionCheckResult2)

		t.FailNow()

	}

}

func Test_CreateTransaction(t *testing.T) {
	fmt.Println("CreateTransaction-----------------")
	NewStub()
	err := MockCreateUser(user1.Id_key, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}
	result1, err := MockCreateTransaction(user1.Id_key, transaction1.Hash, transaction1.Amount, transaction1.Currency_type, transaction1.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction1", result1)

	result2, err := MockCreateTransaction(user1.Id_key, transaction2.Hash, transaction2.Amount, transaction2.Currency_type, transaction2.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}

	fmt.Println("CreateTransaction transaction2", result2)

	result3, err := MockCreateTransaction(user1.Id_key, transaction3.Hash, transaction3.Amount, transaction3.Currency_type, transaction3.Create_at)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction2", result3)

	user, err := MockGetUserAndTransactions(user1.Id_key)
	if err != nil {
		fmt.Println("get User error", err)
	}

	fmt.Println(user)
	assert.Equal(t, len(user.Transactions), 3)

}

func MockCreateUser(Id_key string, name string, email string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateUser"),
			[]byte(Id_key),
			[]byte(name),
			[]byte(email),
		})

	if res.Status != shim.OK {
		fmt.Println("CreateUser failed", string(res.Message))
		return errors.New("CreateUser error")
	}
	return nil
}

func MockGetUserList() ([]*smartcontract.User, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("GetUserList")})
	if res.Status != shim.OK {
		fmt.Println("GetUserList failed", string(res.Message))
		return nil, errors.New("GetUserList error")
	}
	var users []*smartcontract.User
	json.Unmarshal(res.Payload, &users)
	return users, nil
}

func MockUpdateUser(id string, name string, email string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("UpdateUser"),
			[]byte(id),
			[]byte(name),
			[]byte(email),
		})
	if res.Status != shim.OK {
		fmt.Println("UpdateUser failed", string(res.Message))
		return errors.New("UpdateUser error")
	}
	return nil
}

func MockGetUserAndTransactions(id string) (*smartcontract.User, error) {
	var result smartcontract.User
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("GetUserAndTransactions"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("GetUserAndTransactions failed", string(res.Message))
		return nil, errors.New("GetUserAndTransactions error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockCreateTransaction(Id_key string, hash string, amount string, currency_type string, create_at string) (bool, error) {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateTransaction"),
			[]byte(Id_key),
			[]byte(hash),
			[]byte(amount),
			[]byte(currency_type),
			[]byte(create_at),
		})
	if res.Status != shim.OK {
		fmt.Println("CreateTransaction failed", string(res.Message))
		return false, errors.New("CreateTransaction error")
	}
	var result bool = false
	json.Unmarshal(res.Payload, &result)
	return result, nil
}

func MockGetUserByTransactionHash(hash string) (*smartcontract.User, error) {
	var result smartcontract.User
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("TransactionHashExist"),
			[]byte(hash),
		})
	if res.Status != shim.OK {
		fmt.Println("TransactionHashExist failed", string(res.Message))
		return nil, errors.New("TransactionHashExist error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockInitLedger() error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("InitLedger"),
		})
	if res.Status != shim.OK {
		fmt.Println("MockInitLedger failed", string(res.Message))
		return errors.New("MockInitLedger error")
	}
	return nil
}
