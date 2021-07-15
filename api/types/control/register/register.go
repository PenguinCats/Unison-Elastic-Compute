package register

type HandshakeStep1Body struct {
	SecretKey string `json:"secret_key"`
	Seq       int64  `json:"seq"`
}

type HandshakeStep2Body struct {
	Ack   int64  `json:"ack"`
	Seq   int64  `json:"seq"`
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

type HandshakeStep3Body struct {
	Ack int64 `json:"ack"`
}
