package pinksale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchResult struct {
	Docs []struct {
		PoolAddress string `json:"poolAddress"`
	} `json:"docs"`
}

func Search(chainId int64, name string) (result *SearchResult, err error) {
	url := fmt.Sprintf("https://api.pinksale.finance/api/v1/pool/search?chain_id=%d&qs=%s", chainId, name)

	res, err := http.Get(url)
	if err != nil {
		return
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return
}
