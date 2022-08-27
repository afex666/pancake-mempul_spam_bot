package token

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/j0w33x/u-sniper-x/pkg/contracts"
)

type Token struct {
	Name     string
	Symbol   string
	Decimals uint8

	caller     *contracts.ERC20Caller
	transactor *contracts.ERC20Transactor
}

func Load(ether *ethclient.Client, address common.Address) (token *Token, err error) {
	erc20, err := contracts.NewERC20Caller(address, ether)
	if err != nil {
		return
	}

	transactor, err := contracts.NewERC20Transactor(address, ether)
	if err != nil {
		return
	}

	token = new(Token)

	name, err := erc20.Name(nil)
	if err != nil {
		return
	}

	token.Name = name

	symbol, err := erc20.Symbol(nil)
	if err != nil {
		return
	}

	token.Symbol = symbol

	decimals, err := erc20.Decimals(nil)
	if err != nil {
		return
	}

	token.Decimals = decimals

	token.caller = erc20
	token.transactor = transactor

	return
}

func (t *Token) Allowance(from, spender common.Address) (allowance *big.Int, err error) {
	return t.caller.Allowance(nil, from, spender)
}

func (t *Token) BalanceOf(arg0 common.Address) (balance *big.Int, err error) {
	return t.caller.BalanceOf(nil, arg0)
}

func (t *Token) Approve(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int) (tx *types.Transaction, err error) {
	return t.transactor.Approve(opts, arg0, amount)
}
