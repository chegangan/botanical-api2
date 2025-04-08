package e

// GetMsg 根据错误码获取对应的错误消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR_OPERATION_FAILED]
}

// MsgFlags 保存错误码与错误消息的映射关系
var MsgFlags = map[int]string{
	SUCCESS:                       "操作成功",
	ERROR_OPERATION_FAILED:        "操作失败",
	ERROR_INVALID_PARAMS:          "请求参数错误",
	ERROR_USER_NOT_FOUND:          "用户不存在",
	ERROR_USER_ALREADY_EXISTS:     "用户已存在",
	ERROR_USER_CREATE_FAILED:      "创建用户失败",
	ERROR_USER_UPDATE_FAILED:      "更新用户失败",
	ERROR_USER_DELETE_FAILED:      "删除用户失败",
	ERROR_USER_PASSWORD_INCORRECT: "用户名或密码错误",
	ERROR_USER_NO_PERMISSION:      "未授权访问",

	UNAUTHORIZED:                "未授权访问",
	NOT_FOUND:                   "资源不存在",
	INTERNAL_SERVER:             "服务器内部错误",
	ERROR_UPLOAD_FAILED:         "文件上传失败",
	ERROR_PICTURE_NOT_FOUND:     "图片不存在",
	ERROR_PICTURE_NO_PERMISSION: "没有操作权限",
}
