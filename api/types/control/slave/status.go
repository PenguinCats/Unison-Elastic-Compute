/**
 * @File: status
 * @Date: 2021/7/22 下午9:20
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

type StatusSlave int

const (
	SlaveStatusWaitingEstablishControlConnection StatusSlave = iota
	SlaveStatusWaitingEstablishDataConnection
	SlaveStatusNormal
	SlaveStatusOffline
)