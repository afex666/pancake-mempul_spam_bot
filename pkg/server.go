package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	resolver = "http://141.8.199.188"
)

var (
	endpoint      = fmt.Sprintf("%s/v2/validate", resolver)
	fetchEndpoint = fmt.Sprintf("%s/assets/fetch", resolver)
)

type validateLicenseBody struct {
	Input   string `json:"input"`
	Key     string `json:"key"`
	ReqTime int64  `json:"req_time"`
}

type validateLicenseResponse struct {
	Status     bool   `json:"status"`
	Error      string `json:"error,omitempty"`
	Data       string `json:"data,omitempty"`
	ServerTime int64  `json:"server_time,omitempty"`
}

type fetchResponse struct {
	Status  bool   `json:"status"`
	Error   string `json:"error,omitempty"`
	Content string `json:"content,omitempty"`
}

func hash(key, licenseName string, reqTime int64) string {
	hasher := sha256.New()

	hasher.Write([]byte(fmt.Sprintf(
		"%s$$$%s$$$%d",
		key,
		licenseName,
		reqTime,
	)))

	return hex.EncodeToString(hasher.Sum(nil))
}

func hashResponse(key, licenseName string, status bool, serverTime int64) string {
	hasher := sha256.New()

	hasher.Write([]byte(fmt.Sprintf(
		"%s$$$%t$$$%s$$$%d",
		key,
		status,
		licenseName,
		serverTime,
	)))

	return hex.EncodeToString(hasher.Sum(nil))
}

func FetchFile(filename string) (string, error) {
	var resp fetchResponse

	response, err := http.Get(fetchEndpoint)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return "", err
	}

	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	if !resp.Status {
		return "", errors.New("unknown error")
	}

	return resp.Content, nil
}

func ValidSession(key, licenseName string) error {
	var validateResponse validateLicenseResponse

	data := validateLicenseBody{
		Key:     key,
		ReqTime: time.Now().UTC().Unix(),
	}

	data.Input = hash(key, licenseName, data.ReqTime)

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&validateResponse); err != nil {
		return err
	}

	if !validateResponse.Status {
		fmt.Println(validateResponse)
		return errors.New("something bad with server or request")
	}

	if hashResponse(key, licenseName, true, validateResponse.ServerTime) != validateResponse.Data {
		return errors.New("license not found/activated")
	}

	return nil
}
