package yoda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Response struct {
	Text string `json:"text"`
}

func Translate(text string) string {
	vals := url.Values{"text": []string{text}}
	r, err := http.Get(fmt.Sprintf("http://yoda-api.appspot.com/api/v1/yodish?%s", vals.Encode()))
	if err != nil {
		panic(err)
	}

	var res Response
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		panic(err)
	}
	return res.Text
}
