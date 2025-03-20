package e

const (
	Success = 200
	Error   = 500

	// 用户相关错误代码
	ErrorUserNotFound          = 10001
	ErrorUserExists            = 10002
	ErrorUserCreate            = 10003
	ErrorUserUpdate            = 10004
	ErrorUserDelete            = 10005
	ErrorInvalidParams         = 400
	ErrorUserPasswordIncorrect = 10006
	ErrorUserUnauthorized      = 10007
)
