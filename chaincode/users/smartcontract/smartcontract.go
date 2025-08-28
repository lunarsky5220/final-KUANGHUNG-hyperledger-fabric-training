package smartcontract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing archives
type SmartContract struct {
	contractapi.Contract
}

// User Data struct
type User struct {
	Id_key       string        `json:"id_key"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Transactions []Transaction `json:"transactions,omitempty" metadata:",optional"`
}

// Transaction Data struct
type Transaction struct {
	Hash          string `json:"hash"`
	Amount        string `json:"amount"`
	Currency_type string `json:"currency_type"`
	Create_at     string `json:"create_at"`
}

type TransactionHashMapUserId struct {
	Id_key string `json:"Id_key"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id_key string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id_key)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

/*
目前待作的清單
getUserList(完成)
createUser(完成)
getUserAndTransactions(完成)
updateUser(完成)
transactionHashExist(完成)
createTransaction(完成)
*/

func (s *SmartContract) GetUserList(ctx contractapi.TransactionContextInterface) ([]*User, error) {

	//resultsIterator, err := ctx.GetStub().GetStateByRange("0", "99999")
	//按照 字典序 取 [startKey, endKey) 範圍的 key ("","") 代表無邊界取所有值
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		fmt.Println(queryResponse)
		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, id_key string, name string, email string) error {

	exists, err := s.UserExists(ctx, id_key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id_key)
	}

	user := User{
		Name:   name,
		Email:  email,
		Id_key: id_key,
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	//這邊是用id_key當key來存
	return ctx.GetStub().PutState(id_key, userJson)
}

func (s *SmartContract) GetUserAndTransactions(ctx contractapi.TransactionContextInterface, id_key string) (*User, error) {
	userJson, err := ctx.GetStub().GetState(id_key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJson == nil {
		return nil, fmt.Errorf("the user %s does not exist", id_key)
	}

	var user User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id_key string, name string, email string) error {
	user, err := s.GetUserAndTransactions(ctx, id_key)
	if err != nil {
		return err
	}
	user.Email = email
	user.Name = name
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id_key, userJson)
}

func (s *SmartContract) TransactionHashExist(ctx contractapi.TransactionContextInterface, hash string) (*User, error) {
	userJson, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if userJson == nil {
		return nil, fmt.Errorf("the hash %s does not exist", hash)
	}

	var user User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, Id_key string, hash string, amount string, Currency_type string, Create_at string) (bool, error) {
	//撈出指定user
	user, err := s.GetUserAndTransactions(ctx, Id_key)
	if err != nil {
		return false, err
	}

	//將新的transaction 串到這個user上
	var transaction Transaction = Transaction{
		Hash:          hash,
		Amount:        amount,
		Currency_type: Currency_type,
		Create_at:     Create_at,
	}
	user.Transactions = append(user.Transactions, transaction)

	userJson, err := json.Marshal(user)
	if err != nil {
		return false, err
	}
	//轉成json後存回去
	ctx.GetStub().PutState(Id_key, userJson)

	//建立hash->userId的mapping
	var transactionHashMapUserId TransactionHashMapUserId = TransactionHashMapUserId{
		Id_key: user.Id_key,
	}

	transactionHashMapUserIdJson, err := json.Marshal(transactionHashMapUserId)
	if err != nil {
		return false, err
	}

	ctx.GetStub().PutState(hash, transactionHashMapUserIdJson)
	return true, nil
}
