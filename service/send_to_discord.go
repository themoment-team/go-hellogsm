package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"themoment-team/go-hellogsm/internal"
)

const (
	Info  NoticeLevel = "info"
	Warn  NoticeLevel = "warn"
	Error NoticeLevel = "error"

	apiKeyHeader = "x-hg-api-key"
)

type Template struct {
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	NoticeLevel NoticeLevel `json:"noticeLevel"`
}

type NoticeLevel string

func PingRelayApi() error {
	host := internal.SafeApplicationProperties.API.RelayAPI.URL
	resp, err := http.Get(host + "/ping")
	check(err)
	if resp.StatusCode != http.StatusOK {
		return errors.New("relay-api ping failed")
	}
	return nil
}

func SendDiscordMsg(template Template) error {
	host := internal.SafeApplicationProperties.API.RelayAPI.URL
	apikey := internal.SafeApplicationProperties.API.RelayAPI.Key

	marshal, err := json.Marshal(template)
	check(err)
	req, err := http.NewRequest("POST", host+"/notice", bytes.NewReader(marshal))
	check(err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(apiKeyHeader, apikey)

	client := &http.Client{}
	res, err := client.Do(req)
	check(err)

	if res.StatusCode != http.StatusCreated {
		logRequest(marshal)
		return errors.New("relay-api, discord webhook service unavailable")
	}

	return nil
}

func logRequest(marshal []byte) {
	log.Printf("request: %s", string(marshal))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
