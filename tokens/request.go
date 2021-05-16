package tokens

import (
	"io/ioutil"
	"net/http"
)

type accessTokenResponse struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func httpGETRequest(urlStr string) (res []byte, err error) {
	request, err := http.NewRequest(
		"GET",
		urlStr,
		nil,
	)
	if err != nil {

		return
	}

	c := new(http.Client)

	resp, err := c.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	res, err = ioutil.ReadAll(resp.Body)
	return
}
