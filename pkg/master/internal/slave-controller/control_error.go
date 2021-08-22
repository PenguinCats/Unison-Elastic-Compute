/**
 * @File: control_error
 * @Date: 2021/7/15 下午7:40
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import "errors"

var (
	ErrListenerCreat = errors.New("create slave internal_control_types listener fail")

	ErrSlaveControllerCreat = errors.New("create slave controller fail")

	ErrInvalidConnectionRequest = errors.New("invalid internal_connect_types request")

	ErrEstablishCtrlConnInvalidRequest = errors.New("invalid establish ctrl internal_connect_types request")
	ErrEstablishCtrlConnStepFail       = errors.New("establish ctrl internal_connect_types step fail")

	ErrEstablishDataConnInvalidRequest = errors.New("invalid establish internal_data_types internal_connect_types request")
	ErrEstablishDataConnStepFail       = errors.New("establish internal_data_types internal_connect_types step fail")
)
