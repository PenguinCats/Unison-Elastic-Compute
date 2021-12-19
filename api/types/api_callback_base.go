package types

//APICallBackRequestBase 需要回调的 API 请求的 Request base
type APICallBackRequestBase struct {
	OperationID int64  `form:"operation_id" valid:"Required"`
	CallbackURL string `form:"callback_url" valid:"Required"`
}

//APICallBackResponseBase 需要回调的 API 请求的 Response base
type APICallBackResponseBase struct {
	OperationID int64  `json:"operation_id"`
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
}
