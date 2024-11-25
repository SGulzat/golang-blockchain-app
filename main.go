package main

import (
    "github.com/gin-gonic/gin"
    // "gorm.io/driver/mysql"
    "gorm.io/gorm"
	// "github.com/joho/godotenv"
	// "os"
	"fmt"
    "log"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types" 

    "encoding/json"
    "io/ioutil"
    "math/big"
    "strings"
    "context"
)

//---------------------------------------------------------
const (
    rpcURL          = "http://127.0.0.1:8545" // Ganache URL
    privateKeyHex   = "8d068cb5a4d5c74b5c8776a4436ce1e124d23ed5f32e1cc3f01de2b1f58117c9" // Приватный ключ вашего аккаунта без префикса "0x"
    contractAddress = "0x6c30050F1E1D06027887139597d21A36D8FC13B2" // Адрес вашего смарт-контракта
    gasLimit           = uint64(300000)
)


var client *ethclient.Client
var contractAbi abi.ABI
var auth *bind.TransactOpts
var contractAddressHex common.Address
var parsedABI string

//------------------------------------------------------
var db *gorm.DB
var err error

// Task модель для базы данных
type Task struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"` // Например, "open", "in progress", "completed"
    Level      string // "Easy", "Medium", "Hard"
}
//------------------------------------------------------------

type TaskAssignment struct {
    Address common.Address
    ABI     abi.ABI
}
//-------------------------------------------------------
func init() {
    var err error

    // Подключение к клиенту Ethereum
    client, err = ethclient.Dial(rpcURL)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    // Загрузка ABI из файла JSON
    jsonFile, err := ioutil.ReadFile("build/contracts/TaskAssignment.json")
    if err != nil {
        log.Fatalf("Failed to read contract JSON file: %v", err)
    }

    var contractInfo struct {
        ABI      json.RawMessage
        Bytecode string `json:"bytecode"`
    }
    if err := json.Unmarshal(jsonFile, &contractInfo); err != nil {
        log.Fatalf("Failed to unmarshal contract JSON: %v", err)
    }

    // Парсинг ABI
    contractAbi, err = abi.JSON(strings.NewReader(string(contractInfo.ABI)))
    if err != nil {
        log.Fatalf("Failed to parse contract ABI: %v", err)
    }

    // Инициализация аккаунта
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        log.Fatalf("Failed to load private key: %v", err)
    }

    auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // 1337 - это ChainID для Ganache
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }

    // Установка адреса контракта
    contractAddressHex = common.HexToAddress(contractAddress)
}
//-------------------------------------------------------

// ------------------------------------------------------
func main() {
    // // Загружаем переменные окружения из файла .env
    // err := godotenv.Load()
    // if err != nil {
    //     log.Fatal("Ошибка при загрузке файла .env")
    // }

    // // Получаем настройки базы данных из переменных окружения
    // dbUser := os.Getenv("DB_USER")
    // dbPassword := os.Getenv("DB_PASSWORD")
    // dbHost := os.Getenv("DB_HOST")
    // dbPort := os.Getenv("DB_PORT")
    // dbName := os.Getenv("DB_NAME")

    // // Формируем строку подключения
    // dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

    // // Подключение к базе данных
    // db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})		

    // // Миграция модели User
    // db.AutoMigrate(&Task{})

    // // Инициализация роутера
    // r := gin.Default()

    // // CRUD маршруты для Task
    // r.POST("/assign", assignTask)

    // // Запуск сервера
    // r.Run(":8080")
    // Подключение к клиенту Ethereum
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        log.Fatalf("Не удалось подключиться к Ethereum клиенту: %v", err)
    }

    // Загрузка приватного ключа
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        log.Fatalf("Не удалось загрузить приватный ключ: %v", err)
    }

    // Создание авторизованного транзактора
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // Chain ID для основной сети - 1
    if err != nil {
        log.Fatalf("Не удалось создать авторизованного транзактора: %v", err)
    }

    // Загрузка ABI из файла JSON
    jsonFile, err := ioutil.ReadFile("build/contracts/TaskAssignment.json")
    if err != nil {
        log.Fatalf("Failed to read contract JSON file: %v", err)
    }

    var contractInfo struct {
        ABI      json.RawMessage
        Bytecode string `json:"bytecode"`
    }
    if err := json.Unmarshal(jsonFile, &contractInfo); err != nil {
        log.Fatalf("Failed to unmarshal contract JSON: %v", err)
    }

    contractAddr := common.HexToAddress(contractAddress)
    parsedABI, err := abi.JSON(strings.NewReader(string(contractInfo.ABI)))
    if err != nil {
        log.Fatalf("Не удалось распарсить ABI контракта: %v", err)
    }
    _ = contractAddr

    // Пример вызова метода assignTask
    taskID := big.NewInt(101) // Пример ID задачи
    taskLevel := "Easy"

    data, err := parsedABI.Pack("assignTask", taskID, taskLevel)
    if err != nil {
        log.Fatalf("Не удалось упаковать данные для assignTask: %v", err)
    }
    _ = data

        // Добавление разработчика
    developerID := big.NewInt(1) // ID разработчика
    developerLevel := "Junior" // Уровень разработчика

    data2, err := parsedABI.Pack("addDeveloper", developerID, developerLevel)
    if err != nil {
        log.Fatalf("Failed to pack data for addDeveloper: %v", err)
    }

    _ = data2

    //log.Fatalf("Develovers: %v", data2)

    // Получение цены газа
    // gasPrice, err := client.SuggestGasPrice(context.Background())
    // if err != nil {
    //     log.Fatalf("Не удалось получить цену газа: %v", err)
    // }

    // Получение nonce
    nonce, err := client.PendingNonceAt(context.Background(), auth.From)
    if err != nil {
        log.Fatalf("Не удалось получить nonce: %v", err)
    }

    // uint64GasPrice := uint64(1000000000)

    // // Создание транзакции
    // tx := types.NewTransaction(nonce, contractAddr, uint64GasPrice, uint64GasPrice , uint64GasPrice, data)

    // nonce := uint64(1)
    // to := common.HexToAddress("0xRecipientAddress") // Укажите адрес получателя
    value := big.NewInt(1) // 1 ETH в wei big.NewInt(1000000000000000000)
    // gas := uint64(100000) // Стандартное значение газа для простой транзакции
    var gasPrice *big.Int = big.NewInt(20000000000) // 1 Gwei в wei
    gasLimit := uint64(6721975) // Увеличьте до 500000

    // func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction
    tx := types.NewTransaction(nonce, contractAddressHex, value, gasLimit, gasPrice, data)

    // Подписание транзакции
    // signedTx, err := auth.Signer(auth.From, tx)
    // if err != nil {
    //     log.Fatalf("Не удалось подписать транзакцию: %v", err)
    // }

    // // Отправка транзакции
    // err = client.SendTransaction(context.Background(), signedTx)
    // if err != nil {
    //     log.Fatalf("Не удалось отправить транзакцию: %v", err)
    // }

    // Подписка транзакции
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey) // Убедитесь, что вы используете правильный ID цепочки
    if err != nil {
        log.Fatalf("Не удалось подписать транзакцию: %v", err)
    }

    // Отправьте подписанную транзакцию
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        log.Fatalf("Не удалось отправить транзакцию: %v", err)
    }


    fmt.Printf("Транзакция assignTask отправлена: %s\n", signedTx.Hash().Hex())



}


// Приклипление новой задачи
func assignTask(c *gin.Context) {
    // var task Task
    // if err := c.ShouldBindJSON(&task); err != nil {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    //     return
    // }
    // db.Create(&task)
    // contract := bind.NewBoundContract(contractAddressHex, contractAbi, client, client, client)
        
    // txData, err := contractAbi.Pack("assignTask", 1, 2, "task.Level", "Junior") // уровень разработчика можно передавать динамически 

    // if err != nil {
    //     fmt.Errorf("failed to pack data for transaction: %v", err)
    // }

    // // Создание транзакции
    // nonce, err := client.NonceAt(context.Background(), common.HexToAddress(privateKeyHex), nil)
    // if err != nil {
    //     fmt.Errorf("failed to get nonce: %v", err)
    // }

    // // gasPrice, err := client.SuggestGasPrice(context.Background())
    // // if err != nil {
    // //     fmt.Errorf("failed to suggest gas price: %v", err)
    // // }

    // var gasPrice uint64 = 1000000000 

    // tx := types.NewTransaction(nonce, contractAddressHex, big.NewInt(0), gasPrice, nil, txData)

    // c.JSON(http.StatusCreated, "{task:1}")
}
//--------------------------------------------------


