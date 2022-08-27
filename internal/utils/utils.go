package utils

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/j0w33x/u-sniper-x/pkg/token"
)

func StringsToAddresses(in []string) (out []common.Address) {
	for _, element := range in {
		out = append(out, common.HexToAddress(element))
	}

	return
}

func FloatToIntWitDiv(in float64, exp int64) (out *big.Int) {
	inF := big.NewFloat(in)

	outF := new(big.Float).Mul(inF, big.NewFloat(math.Pow(10, float64(exp))))

	out, _ = outF.Int(out)

	return
}

func IntToFloatWithDiv(in *big.Int, exp int64) (out float64) {
	inF := new(big.Float).SetInt(in)

	outF := new(big.Float).Quo(inF, big.NewFloat(math.Pow(10, float64(exp))))

	f, _ := outF.Float64()

	return f
}

func HavePoint(n float64) bool {
	return math.Mod(n, 10.) == 0
}

func Contains(in []int64, arg int64) bool {
	for _, el := range in {
		if el == arg {
			return true
		}
	}

	return false
}

func Reverse(arr []*types.Log) []*types.Log {
	a := make([]*types.Log, 0)

	for i := len(arr); i > 0; i-- {
		a = append(a, arr[i-1])
	}

	return a
}
func ReverseAddr(arr []common.Address) []common.Address {
	a := make([]common.Address, 0)

	for i := len(arr); i > 0; i-- {
		a = append(a, arr[i-1])
	}

	return a
}

func HandleTransactions(hashes []common.Hash, ether *ethclient.Client) (status bool, receipt *types.Receipt) {
	exclude := make([]int64, 0)

	for index, tx := range hashes {
		if len(exclude) >= len(hashes) {
			break
		}

		if Contains(exclude, int64(index)) {
			continue
		}

		receipt, err := ether.TransactionReceipt(context.TODO(), tx)
		for err != nil && err.Error() == "not found" {
			receipt, err = ether.TransactionReceipt(context.TODO(), tx)
		}
		if err != nil && err.Error() != "not found" {
			exclude = append(exclude, int64(index))
			continue
		}

		if receipt.Status == 1 {
			return true, receipt
		}

		exclude = append(exclude, int64(index))
	}

	fmt.Println("Ϡ Не успешно!")
	return false, nil
}

func LiquidityWaitMessage(quit chan *big.Int) *big.Int {
	s := spinner.New(spinner.CharSets[8], 15*time.Millisecond)
	s.Suffix = " Ожидаю ликвидность"
	s.Color("faint", "fgHiCyan")
	s.Start()

	gas := <-quit

	s.Stop()

	fmt.Printf("\033[2K\r؋ Ликвидность обнаружена!\n")

	return gas
}

func SellWaitMessage(amountIn *big.Int, token *token.Token, quit chan bool, amountsGet func() ([]*big.Int, error)) {
	errCount := 0

	amountInF := new(big.Float).SetInt(amountIn)

	s := spinner.New(spinner.CharSets[8], 50*time.Millisecond)
	s.Suffix = " Ожидаю нажатия `q`."
	s.Color("faint", "fgHiCyan")
	s.Start()

	for {
		select {
		case <-quit:
			s.Stop()

			fmt.Println("ω Начинаю продажу!")
			return
		default:
			amounts, err := amountsGet()
			if err != nil && errCount < 10 {
				errCount += 1
				continue
			}

			if errCount >= 10 {
				fmt.Printf("[!] Panic! Err: %s\n", err.Error())
			}

			amountF := new(big.Float).SetInt(amounts[len(amounts)-1])

			amountOut := IntToFloatWithDiv(amounts[len(amounts)-1], int64(token.Decimals))
			profit, _ := new(big.Float).Quo(amountF, amountInF).Float64()

			profit = profit*100 - 100

			s.Suffix = fmt.Sprintf(" Ожидаю нажатия `q`. За ваши токены вы можете получить: %s %s (%s %%)", strconv.FormatFloat(amountOut, 'f', 4, 64), token.Symbol, strconv.FormatFloat(profit, 'f', 2, 32))

			errCount = 0

			time.Sleep(500 * time.Millisecond)
		}
	}
}
