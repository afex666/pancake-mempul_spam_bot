package sniper

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	abiRaw = `[{"inputs":[{"internalType":"address","name":"tokenA","type":"address"},{"internalType":"address","name":"tokenB","type":"address"},{"internalType":"uint256","name":"amountADesired","type":"uint256"},{"internalType":"uint256","name":"amountBDesired","type":"uint256"},{"internalType":"uint256","name":"amountAMin","type":"uint256"},{"internalType":"uint256","name":"amountBMin","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"addLiquidity","outputs":[{"internalType":"uint256","name":"amountA","type":"uint256"},{"internalType":"uint256","name":"amountB","type":"uint256"},{"internalType":"uint256","name":"liquidity","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"amountTokenDesired","type":"uint256"},{"internalType":"uint256","name":"amountTokenMin","type":"uint256"},{"internalType":"uint256","name":"amountETHMin","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"addLiquidityETH","outputs":[{"internalType":"uint256","name":"amountToken","type":"uint256"},{"internalType":"uint256","name":"amountETH","type":"uint256"},{"internalType":"uint256","name":"liquidity","type":"uint256"}],"stateMutability":"payable","type":"function"}]`
	Abi, _ = abi.JSON(strings.NewReader(abiRaw))
)

type decodingResult struct {
	TokenA common.Address
	TokenB common.Address

	AmountA *big.Int
	AmountB *big.Int
}

func decodeInput(weth common.Address, txInput string) (*decodingResult, error) {
	decodedSig, err := hex.DecodeString(txInput[2:10])
	if err != nil {
		return nil, err
	}

	method, err := Abi.MethodById(decodedSig)
	if err != nil {
		return nil, err
	}

	decodedData, err := hex.DecodeString(txInput[10:])
	if err != nil {
		return nil, err
	}

	res, err := method.Inputs.Unpack(decodedData)
	if err != nil {
		return nil, err
	}

	if hexutil.Encode(decodedSig) == "0xe8e33700" {
		tokenA := res[0].(common.Address)
		tokenB := res[1].(common.Address)
		amountA := res[2].(*big.Int)
		amountb := res[3].(*big.Int)

		return &decodingResult{
			TokenA:  tokenA,
			TokenB:  tokenB,
			AmountA: amountA,
			AmountB: amountb,
		}, nil
	}

	if hexutil.Encode(decodedSig) == "0xf305d719" {

		tokenA := weth
		tokenB := res[0].(common.Address)
		amountA := res[3].(*big.Int)
		amountB := res[1].(*big.Int)

		return &decodingResult{
			TokenA:  tokenA,
			TokenB:  tokenB,
			AmountA: amountA,
			AmountB: amountB,
		}, nil
	}

	return nil, nil
}
