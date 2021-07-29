/**
 * @File: message
 * @Date: 2021/7/15 上午10:17
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package control

type MessageCtrlType int

const (
	MessageCtrlTypeHeartbeat MessageCtrlType = iota
	MessageCtrlTypeError
)

type Message struct {
	MessageType MessageCtrlType `json:"message_type"`
	Value       []byte          `json:"value"`
}
