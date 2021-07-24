/**
 * @File: message
 * @Date: 2021/7/15 上午10:17
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package control

type MessageType int

const (
	MessageTypeEstablishCtrlConnection MessageType = iota
	MessageTypeEstablishDataConnection
	MessageTypeReconnect
	MessageTypeError
)

type Message struct {
	MessageType MessageType `json:"message_type"`
	Value       []byte      `json:"value"`
}
