package tokens

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Token is an interface that defines a token
type Token interface {
	renew(cancel context.CancelFunc) error
	ToString() string
}

// auth defines a set of info for authorizing the request
type auth struct {
	id     string
	secret string
}

type accessToken struct {
	auth                auth
	content             string
	expiresAt           int64
	maxConnectionSecond int
}

// NewAccessToken defines the AccessToken for WeCom
func NewAccessToken(corpID, corpSecret string) Token {
	t := new(accessToken)
	t.maxConnectionSecond = 3 // second
	t.auth = auth{
		id:     corpID,
		secret: corpSecret,
	}
	err := t.renew(func() {})
	if err != nil {
		log.Println(err.Error())
	}
	return t
}

func (t *accessToken) renew(cancel context.CancelFunc) error {
	defer cancel()
	res, err := httpGETRequest(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", t.auth.id, t.auth.secret))
	if err != nil {
		return err
	}

	var response accessTokenResponse
	json.Unmarshal(res, &response)
	if response.Errcode != 0 {
		return fmt.Errorf("code: %d, msg: %s", response.Errcode, response.Errmsg)
	}

	t.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second).Unix()
	t.content = response.AccessToken

	return nil
}

// ToString returns the content of token in string
func (t *accessToken) ToString() string {
	if time.Now().After(time.Unix(t.expiresAt-int64(t.maxConnectionSecond), 0)) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.maxConnectionSecond)*time.Second)
		defer cancel()

		go t.renew(cancel)
		<-ctx.Done()
		return t.content
	}

	return t.content
}
