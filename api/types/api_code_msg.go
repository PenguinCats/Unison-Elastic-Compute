package types

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
	UNKNOWN_ERROR  = 999

	SLAVE_INVALID               = 5001
	SLAVE_INSUFFICIENT_RESOURCE = 5002

	CONTAINER_IS_BUSY = 6001
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "内部错误",
	INVALID_PARAMS: "请求参数错误",
	UNKNOWN_ERROR:  "未知错误",

	SLAVE_INVALID:               "设备机不可用",
	SLAVE_INSUFFICIENT_RESOURCE: "设备机资源不足",

	CONTAINER_IS_BUSY: "该容器正在进行其他操作",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
