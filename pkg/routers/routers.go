package routers

import "github.com/ethereum/go-ethereum/common"

type Router struct {
	Name    string
	Address common.Address
}

var (
	RouterList = map[int64][]*Router{
		1: {
			{
				Name:    "UniSwap",
				Address: common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
			},
			{
				Name:    "SushiSwap",
				Address: common.HexToAddress("0xd9e1cE17f2641f24aE83637ab66a2cca9C378B9F"),
			},
		},
		56: {
			{
				Name:    "PancakeSwap",
				Address: common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"),
			},
			{
				Name:    "ApeSwap",
				Address: common.HexToAddress("0xcF0feBd3f17CEf5b47b0cD257aCf6025c5BFf3b7"),
			},
			{
				Name:    "BabySwap",
				Address: common.HexToAddress("0x325e343f1de602396e256b67efd1f61c3a6b38bd"),
			},
			{
				Name:    "BiSwap",
				Address: common.HexToAddress("0x3a6d8ca21d1cf76f653a67577fa0d27453350dd8"),
			},
		},
		137: {
			{
				Name:    "QuickSwap",
				Address: common.HexToAddress("0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff"),
			},
			{
				Name:    "SushiSwap",
				Address: common.HexToAddress("0x1b02dA8Cb0d097eB8D57A175b88c7D8b47997506"),
			},
		},
		43114: {
			{
				Name:    "Trader Joe",
				Address: common.HexToAddress("0x60ae616a2155ee3d9a68541ba4544862310933d4"),
			},
			{
				Name:    "Pangolin",
				Address: common.HexToAddress("0xE54Ca86531e17Ef3616d22Ca28b0D458b6C89106"),
			},
		},
		321: {
			{
				Name:    "QuSwap",
				Address: common.HexToAddress("0xA58350d6dEE8441aa42754346860E3545cc83cdA"),
			},
			{
				Name:    "CoffeeSwap",
				Address: common.HexToAddress("0xc0ffee0000c824d24e0f280f1e4d21152625742b"),
			},
		},
		250: {
			{
				Name:    "SpiritSwap",
				Address: common.HexToAddress("0x16327E3FbDaCA3bcF7E38F5Af2599D2DDc33aE52"),
			},
			{
				Name:    "SpookySwap",
				Address: common.HexToAddress("0xf491e7b69e4244ad4002bc14e878a34207e38c29"),
			},
		},
		25: {
			{
				Name:    "Meerkat Finance",
				Address: common.HexToAddress("0x145677FC4d9b8F19B5D56d1820c48e0443049a30"),
			},
		},
	}
)
