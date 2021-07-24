/**
 * @File: register
 * @Date: 2021/7/24 下午10:02
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package register

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
	UUID  string
	Token string
}

type EstablishDataConnectionHandShakeStep2Body struct {
}
