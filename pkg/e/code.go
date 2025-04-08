package e

// HTTP 标准状态码
const (
	SUCCESS         = 200 // 请求成功
	BAD_REQUEST     = 400 // 请求参数错误
	UNAUTHORIZED    = 401 // 未授权
	FORBIDDEN       = 403 // 禁止访问
	NOT_FOUND       = 404 // 资源不存在
	INTERNAL_SERVER = 500 // 服务器内部错误
)

// 业务错误码 - 用户模块 (1xxxx)
const (
	ERROR_USER_NOT_FOUND          = 10001 // 用户不存在
	ERROR_USER_ALREADY_EXISTS     = 10002 // 用户已存在
	ERROR_USER_CREATE_FAILED      = 10003 // 创建用户失败
	ERROR_USER_UPDATE_FAILED      = 10004 // 更新用户失败
	ERROR_USER_DELETE_FAILED      = 10005 // 删除用户失败
	ERROR_USER_PASSWORD_INCORRECT = 10006 // 用户密码不正确
	ERROR_USER_NO_PERMISSION      = 10007 // 用户无权限
)

// 业务错误码 - 图片模块 (2xxxx)
const (
	ERROR_UPLOAD_FAILED         = 20001 // 上传失败
	ERROR_PICTURE_NOT_FOUND     = 20002 // 图片不存在
	ERROR_PICTURE_NO_PERMISSION = 20003 // 无操作图片权限
)

// 业务错误码 - 系统通用 (9xxxx)
const (
	ERROR_INVALID_PARAMS   = 90001 // 无效的参数
	ERROR_OPERATION_FAILED = 90002 // 操作失败
)
