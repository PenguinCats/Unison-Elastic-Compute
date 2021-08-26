package http_controller

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: 5 * time.Second,
	}
)

func SendCallback(callbackURL string, data []byte) {
	_, err := client.Post(callbackURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logrus.Warningf("send callback message fail with [%s]", err.Error())
	}
}
