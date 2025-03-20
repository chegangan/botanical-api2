package e

// GetMsg 根据错误码获取对应的错误消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}

// MsgFlags 保存错误码与错误消息的映射关系
var MsgFlags = map[int]string{
	Success:                    "操作成功",
	Error:                      "操作失败",
	ErrorInvalidParams:         "请求参数错误",
	ErrorUserNotFound:          "用户不存在",
	ErrorUserExists:            "用户已存在",
	ErrorUserCreate:            "创建用户失败",
	ErrorUserUpdate:            "更新用户失败",
	ErrorUserDelete:            "删除用户失败",
	ErrorUserPasswordIncorrect: "用户名或密码错误",
	ErrorUserUnauthorized:      "未授权访问",
}
