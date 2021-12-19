/**
 * @File: slave_control
 * @Date: 2021/7/15 下午8:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package types

type CreatSlaveBody struct {
	MasterIP        string
	MasterPort      string
	MasterSecretKey string

	HostPortBias int

	Reload bool
}
