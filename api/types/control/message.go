/**
 * @File: message
 * @Date: 2021/7/15 上午10:17
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package control

type MessageCType int

const (
	MessageTypeRegister MessageCType = iota
	MessageTypeReconnect
)

type Message struct {
	MessageType MessageCType `json:"message_type"`
	Value       interface{}  `json:"value"`
}
