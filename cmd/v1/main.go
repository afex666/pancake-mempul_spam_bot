package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/j0w33x/u-sniper-x/internal/config"
	"github.com/j0w33x/u-sniper-x/internal/node"
	"github.com/j0w33x/u-sniper-x/internal/references"
	"github.com/j0w33x/u-sniper-x/internal/utils"
	"github.com/j0w33x/u-sniper-x/pkg"
	"github.com/j0w33x/u-sniper-x/pkg/contracts"
	"github.com/j0w33x/u-sniper-x/pkg/pinksale"
	"github.com/j0w33x/u-sniper-x/pkg/routers"
	"github.com/j0w33x/u-sniper-x/pkg/sniper"
	"github.com/j0w33x/u-sniper-x/pkg/token"
	"github.com/j0w33x/u-sniper-x/pkg/unicrypt"
	"gopkg.in/yaml.v3"
)

func exit(msg string, err error) {
	fmt.Println(msg)
	if err != nil {
		fmt.Printf("৹ Подробней: %s\n", err.Error())
	}

	os.Exit(0)
}

var amountIn *big.Int
var minimalLiquidity *big.Int
var pinkSaleContract = common.HexToAddress("")
var unicryptContract = common.HexToAddress("")

var ABI, _ = abi.JSON(strings.NewReader(contracts.Router02ABI))

func main() {
	defer exit("i Работа завершена!", nil)

	cfg := config.New()
	if err := cfg.Read(references.FILENAME); err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать конфиг(`%s`)!", references.FILENAME), err)
	}

	accounts := config.GetAccounts()

	path := utils.StringsToAddresses(cfg.Settings.Path)
	tokensPath := make([]*token.Token, len(path))

	privateKey, err := crypto.HexToECDSA(cfg.Settings.PrivateKey)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось проверить приватный ключ(`%s`)!", cfg.Settings.PrivateKey), err)
	}

	fmt.Println("೯ Начинаю проверку нод..")

	ether, err := ethclient.DialContext(context.TODO(), cfg.Node.Http)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось подключиться к узлу(`%s`)!", cfg.Node.Http), err)
	}

	rpc, err := rpc.DialContext(context.TODO(), cfg.Node.Ws)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось подключиться к сокету(`%s`)!", cfg.Node.Ws), err)
	}

	if err := node.ValidateNode(ether); err != nil {
		exit(fmt.Sprintf("Hе удалось проверить узел(`%s`)!", cfg.Node.Http), err)
	}

	if err := node.ValidateSocket(rpc); err != nil {
		exit(fmt.Sprintf("Hе удалось проверить сокет(`%s`)!", cfg.Node.Ws), err)
	}

	fmt.Println(" - Успешно проверил ноды!")

	chainId, _ := ether.ChainID(context.TODO())

	chain := node.ByChainID(chainId)
	if chain == nil {
		exit("Неправильная сеть!", errors.New("chain not found"))
		return
	}

	opts, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainId)

	if err := pkg.ValidSession(opts.From.Hex(), "Spamer"); err != nil {
		exit("\nУ вас нет лицензии", errors.New("license not found"))
	}

	fmt.Printf("྾ Сеть: %s (%s)\n", chain.Name, chain.ShortChain)

	if common.HexToAddress(cfg.Settings.Contract).Hex() == common.HexToAddress("").Hex() {
		if _, ok := routers.RouterList[chain.ChainID]; !ok {
			exit("Не удалось найти роутеры под вашу сеть!", err)
		}

		fmt.Printf("\n؋ Адрес смарт-контракта не обнаружен. На вашу сеть есть %d роутера:\n", len(routers.RouterList[chain.ChainID]))
		for index, router := range routers.RouterList[chain.ChainID] {
			fmt.Printf(" %d. %s | %s\n", index+1, router.Name, router.Address.Hex())
		}

		var inputInt int64
		for {
			fmt.Print("Выберите номер > ")
			var input string
			_, err = fmt.Scan(&input)
			if err != nil {
				continue
			}

			inputInt, err = strconv.ParseInt(input, 10, 64)
			if err != nil {
				continue
			}

			if int(inputInt) > len(routers.RouterList[chain.ChainID]) {
				continue
			}

			break
		}

		fmt.Println()

		address, tx, _, err := contracts.DeployFrontCore(opts, ether, routers.RouterList[chain.ChainID][inputInt-1].Address)
		if err != nil {
			exit("Не удалось отправить транзакцию!", err)
		}

		fmt.Printf("ב Отправил транзакцию: %s\n", tx.Hash())

		receipt, err := ether.TransactionReceipt(context.TODO(), tx.Hash())
		for err != nil && err.Error() == "not found" {
			receipt, err = ether.TransactionReceipt(context.TODO(), tx.Hash())
		}

		if err != nil {
			exit("\nНе удалось задеплоить контракт!", err)
		}

		if receipt.Status != 1 {
			exit("\nНе удалось задеплоить контракт!", errors.New("invalid receipt"))
		}

		fmt.Printf("߉ Успешно. Адрес контрактa: %s\n", address)

		cfg.Settings.Contract = address.Hex()

		config, _ := yaml.Marshal(cfg)

		ioutil.WriteFile("config.yaml", config, fs.ModePerm)
	}

	fmt.Println()

	frontCoreCaller, err := contracts.NewFrontCoreCaller(common.HexToAddress(cfg.Settings.Contract), ether)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}

	WETH, err := frontCoreCaller.WETH(nil)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}
	ROUTER, err := frontCoreCaller.RouterAddress(nil)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}

	router, err := contracts.NewRouter02Transactor(ROUTER, ether)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}

	routerCaller, err := contracts.NewRouter02Caller(ROUTER, ether)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}

	frontCore, err := contracts.NewFrontCoreTransactor(common.HexToAddress(cfg.Settings.Contract), ether)
	if err != nil {
		exit(fmt.Sprintf("Hе удалось прочитать контракт(`%s`)!", cfg.Settings.Contract), err)
	}

	fmt.Println("Ξ Путь:")

	for index, tokenAddress := range path {
		tokenInfo, err := token.Load(ether, tokenAddress)
		if err != nil {
			exit(fmt.Sprintf("Hе удалось загрузить токен(`%s`)!", tokenAddress.Hex()), err)
		}

		if index == 0 {
			amountIn = utils.FloatToIntWitDiv(cfg.Settings.Amount, int64(tokenInfo.Decimals))
			minimalLiquidity = utils.FloatToIntWitDiv(cfg.Settings.MinimalLiquidity, int64(tokenInfo.Decimals))
		}
		fmt.Printf(" %d. %s (%s)\n", index+1, tokenInfo.Name, tokenInfo.Symbol)

		tokensPath[index] = tokenInfo
	}

	fmt.Println()

	if inWhiteList, _ := frontCoreCaller.InWhiteList(nil, path[len(path)-1]); !inWhiteList {

		fmt.Println("ζ Добавляю токен в белый лист..")

		approvalTx, err := frontCore.AddWhiteList(opts, path[len(path)-1])
		if err != nil {
			exit(fmt.Sprintf("\nHе удалось апрувнуть токен(`%s`)!", path[len(path)-1].Hex()), err)
		}

		receipt, err := ether.TransactionReceipt(context.TODO(), approvalTx.Hash())
		for err != nil && err.Error() == "not found" {
			receipt, err = ether.TransactionReceipt(context.TODO(), approvalTx.Hash())
		}

		if err != nil {
			exit(fmt.Sprintf("\nHе удалось апрувнуть токен(`%s`)!", path[len(path)-1].Hex()), err)
		}

		if receipt.Status == 1 {
			fmt.Println(" - Успешно!")
		} else {
			fmt.Println(" - Не успешно!")
			return
		}
		fmt.Println()

	}

	opts.GasLimit = cfg.Buy.GasLimit

	amountString := strconv.FormatFloat(cfg.Settings.Amount, 'f', 6, 64)
	amountString = strings.TrimRight(amountString, "0")
	amountString = strings.TrimSuffix(amountString, ".")

	fmt.Printf("ϼ Сумма покупки: %s %s\n", amountString, tokensPath[0].Symbol)

	pinkSaleContractResult, err := pinksale.Search(chain.ChainID, path[len(path)-1].Hex())
	if err != nil {
		exit("pinksale error", err)
	}

	if len(pinkSaleContractResult.Docs) > 0 {
		pinkSaleContract = common.HexToAddress(pinkSaleContractResult.Docs[0].PoolAddress)
	}

	unicryptContractResult, _ := unicrypt.Search(path[len(path)-1].Hex())
	if len(unicryptContractResult.Rows) > 0 {
		unicryptContract = common.HexToAddress(unicryptContractResult.Rows[0].PresaleContract)
	}

	if unicryptContract.Hex() != common.HexToAddress("").Hex() {
		fmt.Println("৳ Найден сейл на UniCrpyt!")
	}

	if pinkSaleContract.Hex() != common.HexToAddress("").Hex() {
		fmt.Println("৳ Найден сейл на PinkSale!")
	}

	fmt.Println()

	quit := make(chan *big.Int, 1)

	listener, _ := sniper.New(cfg.Node.Ws, minimalLiquidity)

	go listener.Start(ether, WETH, ROUTER, path, chainId, pinkSaleContract, unicryptContract, cfg.Signatures, cfg.SigMode, quit)

	hashes := make([]common.Hash, 0)

	gasPrice := utils.LiquidityWaitMessage(quit)

	wg := new(sync.WaitGroup)
	wg.Add(len(accounts))
	for _, account := range accounts {
		go func(wg *sync.WaitGroup, account *ecdsa.PrivateKey, chainId *big.Int) {
			defer wg.Done()

			botOpts, _ := bind.NewKeyedTransactorWithChainID(account, chainId)
			botOpts.GasLimit = cfg.Buy.GasLimit
			botOpts.GasPrice = gasPrice

			for i := 0; i < int(cfg.Settings.TXs); i++ {
				tx, err := frontCore.Swap(botOpts, path, amountIn, minimalLiquidity)
				if err != nil {
					fmt.Println("Не удалось отправить транзакцию!", err)
				}

				hashes = append(hashes, tx.Hash())

				time.Sleep(250 * time.Millisecond)
			}
		}(wg, account, chainId)
	}

	time.Sleep(3 * time.Second)
	wg.Wait()

	fmt.Println()
	fmt.Printf("ζ Отправил %d транзакций. Ожидаю результаты...\n", len(hashes))
	fmt.Println()

	status, receipt := utils.HandleTransactions(hashes, ether)

	if !status {
		tx, err := frontCore.Withdraw(opts, path)
		if err != nil {
			fmt.Println("Не удалось отправить транзакцию на возврат средств! Err: ", err.Error())
			return
		}

		fmt.Printf("\nζ Отправил транзакцию на вывод средств! Hash: %s", tx.Hash())

		return
	}

	logs := utils.Reverse(receipt.Logs)

	var amountOut *big.Int
	var targetToken = tokensPath[len(tokensPath)-1]

	for _, log := range logs {
		if log.Topics[0] != common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822") {
			continue
		}
		swapData, err := ABI.Events["Swap"].Inputs.Unpack(log.Data)
		if err != nil {
			fmt.Println("Не удалось обработать транзакцию! Err: ", err.Error())
			return
		}
		amountOut = swapData[len(swapData)-2].(*big.Int)
	}

	totalBuyedHumanize := utils.IntToFloatWithDiv(amountOut, int64(targetToken.Decimals))

	fmt.Printf("৳ Успешно купил %s %s (%s)\n\n", strconv.FormatFloat(totalBuyedHumanize, 'f', 4, 64), targetToken.Name, targetToken.Symbol)

	if allowance, err := tokensPath[len(tokensPath)-1].Allowance(opts.From, ROUTER); err != nil || amountOut.Cmp(allowance) == 1 {
		opts.GasPrice = nil

		fmt.Println("ζ Апруваю токен продажи")
		approvalTx, err := tokensPath[len(tokensPath)-1].Approve(opts, ROUTER, new(big.Int).Add(amountOut, big.NewInt(1)))
		if err != nil {
			exit(fmt.Sprintf("\nHе удалось апрувнуть токен(`%s`)!", path[len(path)-1].Hex()), err)
		}

		receipt, err = ether.TransactionReceipt(context.TODO(), approvalTx.Hash())
		for err != nil && err.Error() == "not found" {
			receipt, err = ether.TransactionReceipt(context.TODO(), approvalTx.Hash())
		}

		if err != nil {
			exit(fmt.Sprintf("\nHе удалось апрувнуть токен(`%s`)!", path[len(path)-1].Hex()), err)
		}

		if receipt.Status == 1 {
			fmt.Println(" - Успешно!")
		} else {
			fmt.Println(" - Не успешно!")
			fmt.Println()
			return
		}
		fmt.Println()
	}

	reversePath := utils.ReverseAddr(path)

	quitCh := make(chan bool)
	wg = new(sync.WaitGroup)
	wg.Add(1)

	go utils.SellWaitMessage(amountIn, tokensPath[0], quitCh, func() ([]*big.Int, error) {
		return routerCaller.GetAmountsOut(nil, amountOut, reversePath)
	})

	go func() {
		defer wg.Done()

		for {
			ch, r, _ := keyboard.GetSingleKey()

			if ch == 113 && r == 0 {
				quitCh <- true
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)

	wg.Wait()

	opts.GasPrice, _ = ether.SuggestGasPrice(context.TODO())

	opts.GasPrice = new(big.Int).Mul(opts.GasPrice, big.NewInt(2))

	tx, err := router.SwapExactTokensForTokensSupportingFeeOnTransferTokens(opts, amountOut, big.NewInt(0), reversePath, opts.From, big.NewInt(params.Ether))
	if err != nil {
		fmt.Println(" - Не успешно!")
		fmt.Println()
		return
	}

	receipt, err = ether.TransactionReceipt(context.TODO(), tx.Hash())
	for err != nil && err.Error() == "not found" {
		receipt, err = ether.TransactionReceipt(context.TODO(), tx.Hash())
	}

	if err != nil {
		exit(fmt.Sprintf("\nHе удалось продать токен(`%s`)!", path[len(path)-1].Hex()), err)
	}

	if receipt.Status == 1 {
		fmt.Println(" - Успешно!")
	} else {
		fmt.Println(" - Не успешно!")
		fmt.Println()
		return
	}

	fmt.Println()
}
