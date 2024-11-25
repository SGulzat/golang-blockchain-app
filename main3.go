package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "log"
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {
    // Установите соединение с Ganache
    client, err := ethclient.Dial("http://127.0.0.1:8545")
    if err != nil {
        log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
    }

    // Задайте адрес контракта
    contractAddr := common.HexToAddress("0x6c30050F1E1D06027887139597d21A36D8FC13B2") // Ваш адрес контракта

    // Загрузка ABI из файла JSON
    jsonFile, err := ioutil.ReadFile("build/contracts/TaskAssignmentDouble.json")
    if err != nil {
        log.Fatalf("Не удалось прочитать файл контракта: %v", err)
    }

    var contractInfo struct {
        ABI      json.RawMessage
        Bytecode string `json:"bytecode"`
    }
    if err := json.Unmarshal(jsonFile, &contractInfo); err != nil {
        log.Fatalf("Не удалось распарсить JSON контракта: %v", err)
    }

    parsedABI, err := abi.JSON(strings.NewReader(string(contractInfo.ABI)))
    if err != nil {
        log.Fatalf("Не удалось распарсить ABI контракта: %v", err)
    }



    // Параметры для вызова функции assignTask
    taskId := big.NewInt(1)       // ID задачи
    taskLevel := "Easy"            // Уровень задачи
    developerID := big.NewInt(1)   // ID разработчика, которому назначаем задачу

    // Упаковка данных для вызова функции assignTask
    data, err := parsedABI.Pack("assignTask", taskId, taskLevel)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для assignTask: %v", err)
    }

    // Получите nonce
    nonce, err := client.PendingNonceAt(context.Background(), auth.From) // адрес отправителя
    if err != nil {
        log.Fatalf("Не удалось получить nonce: %v", err)
    }

    // Установите лимит газа и цену газа
    gasLimit := uint64(500000) // Лимит газа
    gasPrice := big.NewInt(20000000000) // 20 gwei

    // Создание транзакции
    tx := types.NewTransaction(nonce, contractAddr, nil, gasLimit, gasPrice, data)

    // Получите приватный ключ
    privateKey, err := crypto.HexToECDSA("8d068cb5a4d5c74b5c8776a4436ce1e124d23ed5f32e1cc3f01de2b1f58117c9") // Ваш приватный ключ
    if err != nil {
        log.Fatalf("Не удалось получить приватный ключ: %v", err)
    }

		    // Создайте объект авторизации
    auth := bind.NewKeyedTransactor(privateKey)

    // Подписка транзакции
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey) // Используйте правильный ID сети
    if err != nil {
        log.Fatalf("Не удалось подписать транзакцию: %v", err)
    }

    // Отправьте подписанную транзакцию
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        log.Fatalf("Не удалось отправить транзакцию: %v", err)
    }

    log.Printf("Транзакция отправлена: %s", signedTx.Hash().Hex())

    // Получение результата назначения
    assignedDeveloperID := new(big.Int)
    _, err = client.TransactionReceipt(context.Background(), signedTx.Hash())
    if err != nil {
        log.Fatalf("Не удалось получить результат транзакции: %v", err)
    }

    // Здесь можно добавить логику для получения назначенного разработчика
    // например, вызов getAssignedDeveloper для получения ID назначенного разработчика
    result, err := parsedABI.Pack("getAssignedDeveloper", taskId)
    if err != nil {
        log.Fatalf("Не удалось получить назначенного разработчика: %v", err)
    }
    assignedDeveloperID.Set(result[0].(*big.Int))

    log.Printf("Задача назначена разработчику с ID: %d", assignedDeveloperID)
    // Здесь вы можете сохранить результат в базу данных MySQL
}
