/**
 * @File: register
 * @Date: 2021/7/24 下午10:02
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package internal_connect_types

import "github.com/PenguinCats/Unison-Docker-Controller/api/types/hosts"

type EstablishCtrlConnectionHandshakeStep1Body struct {
	SecretKey string `json:"secret_key"`
	Seq       int64  `json:"seq"`
}

type EstablishCtrlConnectionHandshakeStep2Body struct {
	Ack   int64  `json:"ack"`
	Seq   int64  `json:"seq"`
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

type EstablishCtrlConnectionHandshakeStep3Body struct {
	Ack int64 `json:"ack"`
}

type EstablishDataConnectionHandShakeStep1Body struct {
	UUID     string
	Token    string
	HostInfo hosts.HostInfo
}

type EstablishDataConnectionHandShakeStep2Body struct {
}

type ReconnectCtrlConnectionHandshakeStep1Body struct {
	Seq   int64  `json:"seq"`
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

type ReconnectCtrlConnectionHandshakeStep2Body struct {
	Ack   int64 `json:"ack"`
	Agree bool  `json:"agree"`
}
