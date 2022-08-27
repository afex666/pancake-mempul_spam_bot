package unicrypt

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SearchResult struct {
	Count string `json:"count"`
	Rows  []struct {
		PresaleContract string `json:"presale_contract"`
	} `json:"rows"`
}

func Search(name string) (result *SearchResult, err error) {
	url := "https://api-pcakev2.unicrypt.network/api/v1/presales/search"

	body := map[string]interface{}{
		"filters": map[string]interface{}{
			"hide_flagged":  false,
			"search":        name,
			"show_hidden":   false,
			"sort":          "uncl_participants",
			"sortAscending": true,
		},
		"page":          0,
		"rows_per_page": 6,
		"stage":         0,
	}

	jsonBody, _ := json.Marshal(body)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return
}
