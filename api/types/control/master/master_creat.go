/**
 * @File: slave_control
 * @Date: 2021/7/15 下午8:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package master

type CreatMasterBody struct {
	SlaveControlListenerPort string

	MasterAPIPort string
}
