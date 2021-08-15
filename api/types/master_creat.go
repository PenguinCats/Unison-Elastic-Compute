/**
 * @File: slave_control
 * @Date: 2021/7/15 下午8:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package types

type CreatMasterBody struct {
	Recovery string

	SlaveControlListenerPort string

	APIPort string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}
