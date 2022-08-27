package node

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	validateTime = time.Duration(2) * time.Second
)

func ValidateNode(ether *ethclient.Client) (err error) {
	if _, err = ether.HeaderByNumber(context.TODO(), nil); err != nil {
		return
	}

	return
}

func ValidateSocket(rpc *rpc.Client) (err error) {
	pendingTxs := make(chan common.Hash)
	_, err = rpc.Subscribe(context.TODO(), "eth", pendingTxs, "newPendingTransactions")
	if err != nil {
		return
	}

	isValid := make(chan bool)

	txs := make([]common.Hash, 0)

	go func(isValid chan bool, pendingChan chan common.Hash, txs *[]common.Hash) {
		for {
			select {
			case <-isValid:
				return
			case tx := <-pendingChan:
				*txs = append(*txs, tx)
				isValid <- true
			}
		}
	}(isValid, pendingTxs, &txs)

	time.Sleep(validateTime)

	if len(txs) == 0 {
		return fmt.Errorf("0 pending txs for %f secs", validateTime.Seconds())
	}

	return nil
}
