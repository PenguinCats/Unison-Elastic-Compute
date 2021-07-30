/**
 * @File: connect
 * @Date: 2021/7/24 下午9:31
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package connect

type ConnectionType int

const (
	ConnectionTypeEstablishCtrlConnection ConnectionType = iota
	ConnectionTypeEstablishDataConnection
	ConnectionTypeReconnect
	ConnectionTypeError
)

type ConnectionHead struct {
	ConnectionType ConnectionType `json:"connection_type"`
}
