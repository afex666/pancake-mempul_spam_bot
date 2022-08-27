package config

import (
	"crypto/ecdsa"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Node struct {
		Http string `yaml:"http"`
		Ws   string `yaml:"ws"`
	} `yaml:"node"`
	Settings struct {
		PrivateKey       string   `yaml:"privateKey"`
		MinimalLiquidity float64  `yaml:"minimalLiquidity"`
		Path             []string `yaml:"path"`
		Amount           float64  `yaml:"amount"`
		Contract         string   `yaml:"contract"`

		TXs int64 `yaml:"txs"`
	} `yaml:"settings"`
	Buy struct {
		GasLimit uint64 `yaml:"gasLimit"`
	} `yaml:"buy"`
	SigMode    bool     `yaml:"sigmode"`
	Signatures []string `yaml:"signatures"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Read(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, &c)
}

func GetAccounts() (privateKeys []*ecdsa.PrivateKey) {
	privateKeys = make([]*ecdsa.PrivateKey, 0)

	content, err := ioutil.ReadFile("accounts.txt")
	if err != nil {
		return
	}

	data := strings.Split(string(content), "\n")
	if len(data) == 0 {
		return
	}

	for _, line := range data {
		if len(line) < 42 {
			continue
		}

		line = strings.TrimPrefix(line, "0x")
		line = strings.ReplaceAll(line, "\n", "")
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\t", "")

		privateKey, err := crypto.HexToECDSA(line)
		if err != nil {
			continue
		}

		privateKeys = append(privateKeys, privateKey)
	}

	return
}
