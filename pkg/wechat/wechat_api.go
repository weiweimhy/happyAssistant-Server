package wechat

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type WeChatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

const URL_FORMAT = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

func GetSession(appid, secret, code string, callback func(*WeChatSession)) error {
	url := fmt.Sprintf(URL_FORMAT, appid, secret, code)
	rsp, err := http.Get(url)
	if err != nil {
		return err
	}

	if rsp.Body != nil {
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				log.Errorf("Get wechat session close body error: %v", err)
			}
		}(rsp.Body)
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("Get wechat session status code: %d", rsp.StatusCode)
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	session := &WeChatSession{}
	err = json.Unmarshal(data, session)
	if err != nil {
		return err
	}
	callback(session)

	return nil
}
