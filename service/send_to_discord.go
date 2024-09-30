package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"themoment-team/go-hellogsm/internal"
)

const (
	Info  NoticeLevel = "info"
	Warn  NoticeLevel = "warn"
	Error NoticeLevel = "error"

	EnvDev  Env = "dev"
	EnvProd Env = "prod"

	apiKeyHeader = "x-hg-api-key"
)

type Template struct {
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	NoticeLevel NoticeLevel `json:"noticeLevel"`
	Env         Env         `json:"env"`
}

type NoticeLevel string
type Env string

func GetEnv() Env {
	var env Env
	if internal.GetActiveProfile().Value == internal.Local.Value ||
		internal.GetActiveProfile().Value == internal.Stage.Value {
		env = EnvDev
	} else if internal.GetActiveProfile().Value == internal.Prod.Value {
		env = EnvProd
	} else {
		panic("현재 설정되어있는 profile이 정상적이지 않습니다.")
	}
	return env
}

// PingRelayApi
// 더모먼트팀 RelayAPI 로 Ping 요청을 보낸다.
func PingRelayApi() error {
	host := internal.SafeApplicationProperties.API.RelayAPI.URL
	resp, err := http.Get(host + "/ping")
	check(err)
	if resp.StatusCode != http.StatusOK {
		return errors.New("relay-api ping failed")
	}
	return nil
}

// SendDiscordMsg
// 더모먼트팀 RelayAPI 를 통해 discord 메시지를 보낸다.
func SendDiscordMsg(template Template) {
	host := internal.SafeApplicationProperties.API.RelayAPI.URL
	apikey := internal.SafeApplicationProperties.API.RelayAPI.Key

	standardization(&template)

	apiClientErr := postMessageToDiscord(template, host, apikey)
	if apiClientErr != nil {
		log.Printf("[%s] <- 메시지 전송 처리 중 오류 발생. 상세: [%s]", template.Content, apiClientErr.Error())
	}
}

// 디스코드로 전송할 메시지를 규격에 맞게 변경한다.
func standardization(template *Template) {
	template.Title = addTitlePrefix(template.Title)
}

func addTitlePrefix(title string) string {
	return fmt.Sprintf("[%s] %s", internal.GetActiveProfile().Desc, title)
}

func postMessageToDiscord(template Template, host string, apikey string) error {
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
