package main

import (
    "context"
    "crypto/ecdsa"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum"
)

func main() {
    // Подключение к Ganache
    client, err := ethclient.Dial("http://127.0.0.1:8545")
    if err != nil {
        log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
    }

    // Адрес контракта
    contractAddr := common.HexToAddress("0xAa706a75c8b0eA294cf300729099eD231372fEc7") // Замените на ваш адрес контракта

    // Загрузка ABI контракта
    contractJSON, err := ioutil.ReadFile("build/contracts/TaskAssignmentPro.json")
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

    // Приватный ключ
    privateKey, err := crypto.HexToECDSA("8f733ed832243f4bb26720822ab0a682afbee6e831ce435cc859f3544e6ffef0") // Замените на ваш приватный ключ
    if err != nil {
        log.Fatalf("Не удалось получить приватный ключ: %v", err)
    }

    // Адрес отправителя
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("Не удалось преобразовать публичный ключ в ECDSA")
    }
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

    // Получение nonce
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        log.Fatalf("Не удалось получить nonce: %v", err)
    }

    // Параметры транзакции
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatalf("Не удалось получить ID сети: %v", err)
    }

    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatalf("Не удалось получить цену газа: %v", err)
    }

	// Соответствие значений enum:
	// DeveloperLevel:
	// 0 - Junior
	// 1 - Middle
	// 2 - Senior

	developers := []struct {
		id    *big.Int
		level uint8
	}{
		{big.NewInt(1), 0}, // Junior
		{big.NewInt(2), 1}, // Middle
		{big.NewInt(3), 2}, // Senior
	}


    for _, developer := range developers {
        // Упаковка данных для вызова функции addDeveloper
        input, err := parsedABI.Pack("addDeveloper", developer.id, developer.level)
        if err != nil {
            log.Fatalf("Не удалось упаковать данные для addDeveloper: %v", err)
        }

        // Создание транзакции
        tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(500000), gasPrice, input)

        // Подписание транзакции
        signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
        if err != nil {
            log.Fatalf("Не удалось подписать транзакцию: %v", err)
        }

        // Отправка транзакции
        err = client.SendTransaction(context.Background(), signedTx)
        if err != nil {
            log.Fatalf("Не удалось отправить транзакцию: %v", err)
        }

        fmt.Printf("Разработчик %d добавлен, транзакция: %s\n", developer.id, signedTx.Hash().Hex())

        nonce++ // Увеличиваем nonce для следующей транзакции
    }

	// Соответствие значений enum:
	// TaskLevel:
	// 0 - Easy
	// 1 - Medium
	// 2 - Hard

	taskID := big.NewInt(1)
	taskLevel := uint8(1) // Easy

    // Упаковка данных для вызова функции assignTask
    input, err := parsedABI.Pack("assignTask", taskID, taskLevel)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для assignTask: %v", err)
    }

    // Создание транзакции
    tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(500000), gasPrice, input)

    // Подписание транзакции
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        log.Fatalf("Не удалось подписать транзакцию: %v", err)
    }

    // Отправка транзакции
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        log.Fatalf("Не удалось отправить транзакцию: %v", err)
    }

    fmt.Printf("Задача %d назначена, транзакция: %s\n", taskID, signedTx.Hash().Hex())

    nonce++

    // Вызов функции getAssignedDeveloper
    input, err = parsedABI.Pack("getAssignedDeveloper", taskID)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для getAssignedDeveloper: %v", err)
    }

    // Подготовка сообщения для вызова контракта
    callMsg := ethereum.CallMsg{
        From: fromAddress,
        To:   &contractAddr,
        Data: input,
    }

		// Выполнение вызова
	output, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Не удалось выполнить вызов контракта: %v", err)
	}

	// Распаковка результата
	results, err := parsedABI.Unpack("getAssignedDeveloper", output)
	if err != nil {
		log.Fatalf("Не удалось распаковать данные: %v", err)
	}

	if len(results) > 0 {
		assignedDeveloperID := results[0].(*big.Int)
		fmt.Printf("Задача %d назначена разработчику с ID: %d\n", taskID, assignedDeveloperID)
	} else {
		log.Println("Не удалось получить ID назначенного разработчика")
	}
}