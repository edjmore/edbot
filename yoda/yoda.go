package yoda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Response struct {
	Yodish string `json:"yodish"`
}

func Translate(text string) string {
	vals := url.Values{"text": []string{text}}
	r, err := http.Get(fmt.Sprintf("http://yoda-api.appspot.com/api/v1/yodish?%s", vals.Encode()))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", r)

	var res Response
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		panic(err)
	}
	return res.Yodish
}
