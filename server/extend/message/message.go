package message

import "net/http"

// Code 错误输出数据结构
type Code struct {
	Status    int    `json:"status"`  // HTTP 状态
	Code      int    `json:"code"`    // 业务错误码
	MessageEn string `json:"message"` // 业务错误信息
	MessageCn string `json:"text"`    // 业务错误信息
}

var (
	Success     = &Code{http.StatusOK, 0, "success", "成功"}
	Failed      = &Code{http.StatusOK, -1, "failed", "失败"}
	ServerError = &Code{http.StatusOK, 500, "ServerError", "服务器内部异常"}
	// 参数相关（200~299）
	ParamErr   = &Code{http.StatusOK, 200, "param err", "参数错误"}
	LoginError = &Code{http.StatusOK, 302, "account or password error", "用户名或者密码错误"}
	// 参数相关（200~299）
	RequestFail = &Code{http.StatusOK, 402, "request fail", "请求失败，请检查"}
)
