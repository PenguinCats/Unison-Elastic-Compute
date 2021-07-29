/**
 * @File: control_error
 * @Date: 2021/7/15 下午7:40
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import "errors"

var (
	ErrListenerCreat = errors.New("create slave control listener fail")

	ErrSlaveControllerCreat = errors.New("create slave controller fail")

	ErrInvalidConnectionRequest = errors.New("invalid connect request")

	ErrEstablishCtrlConnInvalidRequest = errors.New("invalid establish ctrl connect request")
	ErrEstablishCtrlConnStepFail       = errors.New("establish ctrl connect step fail")

	ErrEstablishDataConnInvalidRequest = errors.New("invalid establish data connect request")
	ErrEstablishDataConnStepFail       = errors.New("establish data connect step fail")

	ErrControlInvalidMessage   = errors.New("invalid control message")
	ErrControlInvalidHeartbeat = errors.New("invalid heartbeat message")
)
