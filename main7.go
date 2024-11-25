package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/joho/godotenv"
	"os"

)

type TaskAssignmentRequest struct {
	Developers []struct {
		ID    int64 `json:"id"`
		Level uint8 `json:"level"`
	} `json:"developers"`
	TaskID    int64 `json:"task_id"`
	TaskLevel uint8 `json:"task_level"`
}

// Функция для назначения задачи с использованием данных из запроса
func assignTaskHandler(c *gin.Context) {
	var req TaskAssignmentRequest

	// Парсинг JSON-запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Загружаем переменные окружения из файла .env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке файла .env")
    }

    // Получаем настройки базы данных из переменных окружения
    smartContractAdress := os.Getenv("CONTRACT_ADRESS")
    privateKeyAdress := os.Getenv("PRIVATE_KEY_HEX")
	rpcURL := os.Getenv("RPC_URL")

    // Формируем строку подключения
    fmt.Printf("Private Key: %s\n", privateKeyAdress)
	fmt.Printf("Contract Address: %s\n", smartContractAdress)

    client, err := ethclient.Dial(rpcURL)
    _ = client
	if err != nil {
		log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
	}

	contractAddr := common.HexToAddress(smartContractAdress)
    _ = contractAddr

	// Загрузка ABI контракта
	contractJSON, err := ioutil.ReadFile("build/contracts/TaskAssignmentProMax.json")
    _ = contractJSON
	if err != nil {
		log.Fatalf("Не удалось прочитать файл контракта: %v", err)
	}

    var contractInfo struct {
		ABI json.RawMessage
	}

	err = json.Unmarshal(contractJSON, &contractInfo)
	if err != nil {
		log.Fatalf("Не удалось распарсить JSON контракта: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractInfo.ABI)))
    _ = parsedABI
	if err != nil {
		log.Fatalf("Не удалось распарсить ABI контракта: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyAdress)
    _ = privateKey
	if err != nil {
		log.Fatalf("Не удалось получить приватный ключ: %v", err)
	}

	/*
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
	}

	contractAddr := common.HexToAddress("0x9C2a54127273b7de3694Fe0254190A37ed586c5a")

	// Загрузка ABI контракта
	contractJSON, err := ioutil.ReadFile("build/contracts/TaskAssignmentProMax.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл контракта: %v", err)
	}

	var contractInfo struct {
		ABI json.RawMessage
	}

	err = json.Unmarshal(contractJSON, &contractInfo)
	if err != nil {
		log.Fatalf("Не удалось распарсить JSON контракта: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractInfo.ABI)))
	if err != nil {
		log.Fatalf("Не удалось распарсить ABI контракта: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("f0a80a4c991547dd733661ecb8aaff109fadf96f250d2d99796a106870e11928")
	if err != nil {
		log.Fatalf("Не удалось получить приватный ключ: %v", err)
	}
	*/

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Не удалось преобразовать публичный ключ в ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Не удалось получить nonce: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Не удалось получить ID сети: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Не удалось получить цену газа: %v", err)
	}

	// **Вызов функции resetDevelopers()**
	input, err := parsedABI.Pack("resetDevelopers")
	if err != nil {
		log.Fatalf("Не удалось упаковать данные для resetDevelopers: %v", err)
	}

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(500000), gasPrice, input)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Не удалось подписать транзакцию: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Не удалось отправить транзакцию: %v", err)
	}

	fmt.Printf("Список разработчиков сброшен, транзакция: %s\n", signedTx.Hash().Hex())

	// **Дожидаемся подтверждения транзакции resetDevelopers() (опционально)**
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatalf("Ошибка при ожидании подтверждения транзакции: %v", err)
	}
	if receipt.Status != 1 {
		log.Fatalf("Транзакция resetDevelopers() не выполнена успешно")
	}

	nonce++ // Увеличиваем nonce после отправки транзакции

	// **Добавление разработчиков**
	for _, developer := range req.Developers {
		input, err := parsedABI.Pack("addDeveloper", big.NewInt(developer.ID), developer.Level)
		if err != nil {
			log.Fatalf("Не удалось упаковать данные для addDeveloper: %v", err)
		}

		tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(500000), gasPrice, input)

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatalf("Не удалось подписать транзакцию: %v", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatalf("Не удалось отправить транзакцию: %v", err)
		}

		fmt.Printf("Разработчик %d добавлен, транзакция: %s\n", developer.ID, signedTx.Hash().Hex())
		nonce++
	}

	// **Дожидаемся подтверждения транзакций addDeveloper (опционально)**
	// Вы можете добавить ожидание подтверждения транзакций для каждого разработчика, если это необходимо.

	// **Назначение задачи**
	taskID := big.NewInt(req.TaskID)
	taskLevel := req.TaskLevel

	input, err = parsedABI.Pack("assignTask", taskID, taskLevel)
	if err != nil {
		log.Fatalf("Не удалось упаковать данные для assignTask: %v", err)
	}

	tx = types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(500000), gasPrice, input)

	signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Не удалось подписать транзакцию: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Не удалось отправить транзакцию: %v", err)
	}

	fmt.Printf("Задача %d назначена, транзакция: %s\n", req.TaskID, signedTx.Hash().Hex())

	// **Дожидаемся подтверждения транзакции assignTask (опционально)**
	receipt, err = bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatalf("Ошибка при ожидании подтверждения транзакции: %v", err)
	}
	if receipt.Status != 1 {
		log.Fatalf("Транзакция assignTask не выполнена успешно")
	}

	nonce++ // Увеличиваем nonce после отправки транзакции

	// **Получение назначенного разработчика**
	input, err = parsedABI.Pack("getAssignedDeveloper", taskID)
	if err != nil {
		log.Fatalf("Не удалось упаковать данные для getAssignedDeveloper: %v", err)
	}

	callMsg := ethereum.CallMsg{
		From: fromAddress,
		To:   &contractAddr,
		Data: input,
	}

	output, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Не удалось выполнить вызов контракта: %v", err)
	}

	results, err := parsedABI.Unpack("getAssignedDeveloper", output)
	if err != nil {
		log.Fatalf("Не удалось распаковать данные: %v", err)
	}

	if len(results) > 0 {
		assignedDeveloperID := results[0].(*big.Int)
		fmt.Printf("Задача %d назначена разработчику с ID: %d\n", req.TaskID, assignedDeveloperID)
		c.JSON(http.StatusOK, gin.H{"message": "Task assigned", "developer_id": assignedDeveloperID.String()})
	} else {
		log.Println("Не удалось получить ID назначенного разработчика")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assigned developer"})
	}
}


func main() {
	r := gin.Default()
	r.POST("/assign-task", assignTaskHandler) // Используем POST для получения данных в JSON формате
	r.Run(":8080")
}