package main

import (
    "context"
    "crypto/ecdsa"
    "encoding/json"
    "io/ioutil"
    "log"
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/keystore"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/rpc"
)

func main() {
    // Переменные для подключения к Ethereum
    var (
        client         *ethclient.Client
        err            error
        privateKey     *ecdsa.PrivateKey
        contractAddr    common.Address
        jsonFile       []byte
        contractInfo    struct {
            ABI      json.RawMessage
            Bytecode string `json:"bytecode"`
        }
        parsedABI      abi.ABI
        auth           *bind.TransactOpts
    )

    // Переменные для разработчиков
    var (
        developers = []struct {
            id    *big.Int
            level string
        }{
            {big.NewInt(1), "Junior"},
            {big.NewInt(2), "Middle"},
            {big.NewInt(3), "Senior"},
        }
    )

    // Переменные для транзакций
    var (
        nonce     uint64
        gasLimit  uint64 = 500000      // Лимит газа для транзакций
        gasPrice  = big.NewInt(20000000000) // Цена газа в wei (20 gwei)
        tx        *types.Transaction
        tx2        *types.Transaction
        signedTx  *types.Transaction
        signedTx2  *types.Transaction
        result    []byte
        assignedDeveloperID *big.Int
    )

    // Переменные для задач
    var (
        taskId    = big.NewInt(1)  // ID задачи
        taskLevel = "Easy"          // Уровень задачи
        receipt   *types.Receipt     // Результат транзакции
    )

    // 1. Подключение к Ganache
    client, err = ethclient.Dial("http://127.0.0.1:7545")
    if err != nil {
        log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
    }

    // 2. Установите адрес контракта
    contractAddr = common.HexToAddress("0x6c30050F1E1D06027887139597d21A36D8FC13B2") // Ваш адрес контракта

    // 3. Загрузка ABI из файла JSON
    jsonFile, err = ioutil.ReadFile("build/contracts/TaskAssignmentDouble.json")
    if err != nil {
        log.Fatalf("Не удалось прочитать файл контракта: %v", err)
    }

    // 4. Распаковка ABI
    if err := json.Unmarshal(jsonFile, &contractInfo); err != nil {
        log.Fatalf("Не удалось распарсить JSON контракта: %v", err)
    }

    parsedABI, err = abi.JSON(strings.NewReader(string(contractInfo.ABI)))
    if err != nil {
        log.Fatalf("Не удалось распарсить ABI контракта: %v", err)
    }

    // 5. Получите приватный ключ
    privateKey, err = crypto.HexToECDSA("8d068cb5a4d5c74b5c8776a4436ce1e124d23ed5f32e1cc3f01de2b1f58117c9") // Ваш приватный ключ
    if err != nil {
        log.Fatalf("Не удалось получить приватный ключ: %v", err)
    }

    // 6. Создайте объект авторизации
    auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // Используйте правильный ID сети
    if err != nil {
        log.Fatalf("Не удалось создать авторизацию: %v", err)
    }

    for _, developer := range developers {
        // Упаковка данных для вызова функции addDeveloper
        addDevData, err := parsedABI.Pack("addDeveloper", developer.id, developer.level)
        if err != nil {
            log.Fatalf("Не удалось упаковать данные для addDeveloper: %v", err)
        }

        // Получите nonce
        nonce, err := client.PendingNonceAt(context.Background(), auth.From) // адрес отправителя
        if err != nil {
            log.Fatalf("Не удалось получить nonce: %v", err)
        }

        // Установите лимит газа и цену газа
        gasLimit := uint64(500000) // Лимит газа
        gasPrice := big.NewInt(20000000000) // 20 gwei

        // Создание транзакции для добавления разработчика
        tx := types.NewTransaction(nonce, contractAddr, nil, gasLimit, gasPrice, addDevData)

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

        log.Printf("Транзакция для добавления разработчика отправлена: %s", signedTx.Hash().Hex())
    }

    // 2. Назначение задачи
    // taskId := big.NewInt(1)       // ID задачи
    // taskLevel := "Easy"            // Уровень задачи

    // Упаковка данных для вызова функции assignTask
    assignTaskData, err := parsedABI.Pack("assignTask", taskId, taskLevel)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для assignTask: %v", err)
    }

    // Получите nonce
    nonce, err = client.PendingNonceAt(context.Background(), auth.From) // адрес отправителя
    if err != nil {
        log.Fatalf("Не удалось получить nonce: %v", err)
    }

    // Создание транзакции для назначения задачи
    tx2 := types.NewTransaction(nonce, contractAddr, nil, gasLimit, gasPrice, assignTaskData)

    // Подписка транзакции
    signedTx2, err = types.SignTx(tx2, types.NewEIP155Signer(big.NewInt(1)), privateKey) // Используйте правильный ID сети
    if err != nil {
        log.Fatalf("Не удалось подписать транзакцию: %v", err)
    }

    // Отправьте подписанную транзакцию
    err = client.SendTransaction(context.Background(), signedTx2)
    if err != nil {
        log.Fatalf("Не удалось отправить транзакцию: %v", err)
    }

    log.Printf("Транзакция назначения задачи отправлена: %s", signedTx2.Hash().Hex())

    // Получение результата назначения
    receipt, err := client.TransactionReceipt(context.Background(), signedTx2.Hash())
    if err != nil {
        log.Fatalf("Не удалось получить результат транзакции: %v", err)
    }

    // Здесь можно проверить, была ли задача успешно назначена
    if receipt.Status == 1 {
        log.Println("Задача успешно назначена.")
    } else {
        log.Println("Не удалось назначить задачу.")
    }

    // Получение ID назначенного разработчика
    assignedDeveloperID := new(big.Int)
    // Упаковка данных для вызова функции getAssignedDeveloper
    getAssignedData, err := parsedABI.Pack("getAssignedDeveloper", taskId)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для getAssignedDeveloper: %v", err)
    }

    // Создание вызова
    callMsg := ethereum.CallMsg{
        To:   &contractAddr,
        Data: getAssignedData,
    }

    // Выполнение вызова
    result, err := client.CallContract(context.Background(), callMsg, nil)
    if err != nil {
        log.Fatalf("Не удалось выполнить вызов контракта: %v", err)
    }

    // Распаковка результата
    err = parsedABI.UnpackIntoInterface(assignedDeveloperID, "getAssignedDeveloper", result)
    if err != nil {
        log.Fatalf("Не удалось распаковать результат: %v", err)
    }

    log.Printf("Задача назначена разработчику с ID: %d", assignedDeveloperID)
    // Здесь вы можете сохранить результат в базу данных MySQL
}
