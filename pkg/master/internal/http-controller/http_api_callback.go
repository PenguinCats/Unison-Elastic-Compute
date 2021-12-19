package http_controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: 5 * time.Second,
	}
)

//func SendCallback(callbackURL string, data []byte) {
//	_, err := client.Post(callbackURL, "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		logrus.Warningf("send callback message fail with [%s]", err.Error())
//	}
//}

func SendCallbackPostWithoutResponse(callbackURL string, payload interface{}) error {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(callbackURL, "application/json", bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("http error")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return nil
}
