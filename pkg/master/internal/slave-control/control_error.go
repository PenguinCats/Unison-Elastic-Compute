/**
 * @File: control_error
 * @Date: 2021/7/15 下午7:40
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import "errors"

var (
	ErrSlaveControllerCreat = errors.New("create slave controller fail")

	ErrListenerCreat = errors.New("create slave control listener fail")

	ErrConnectionRequestInvalid      = errors.New("connect request is invalid")
	ErrConnectionRequestWrongAddress = errors.New("connect request has a wrong address")

	ErrRegisterInvalidBody = errors.New("register body is invalid")
	ErrRegisterInvalidInfo = errors.New("register info is invalid")
)
