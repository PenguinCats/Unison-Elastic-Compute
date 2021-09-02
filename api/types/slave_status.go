/**
 * @File: status
 * @Date: 2021/7/22 下午9:20
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package types

type StatsSlave int

const (
	StatsWaitingEstablishControlConnection StatsSlave = iota
	StatsWaitingEstablishDataConnection
	StatsNormal
	StatsStopped
	StatsOffline
)

func GetSlaveStatsString(stats StatsSlave) string {
	switch stats {
	case StatsWaitingEstablishControlConnection, StatsWaitingEstablishDataConnection:
		return "connecting"
	case StatsNormal:
		return "online"
	case StatsStopped:
		return "stopped"
	case StatsOffline:
		return "offline"
	default:
		return "error"
	}
}
