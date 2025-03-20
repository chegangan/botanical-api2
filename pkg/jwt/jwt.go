package jwt

import (
	"botanical-api2/pkg/setting"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// 定义标准错误
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("无效的令牌")
)

// CustomClaims 自定义JWT声明结构体
type CustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int, username string) (string, error) {
	// 从配置文件获取过期时间
	expireTime := time.Now().Add(time.Duration(setting.JwtExpireHours) * time.Hour)

	claims := CustomClaims{
		ID:       userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(), // 签发时间
			Issuer:    "botanical-api",   // 签发人
		},
	}

	// 创建令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌
	token, err := tokenClaims.SignedString([]byte(setting.JwtSecret))

	return token, err
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(setting.JwtSecret), nil
	})

	if err != nil {
		// 处理常见的错误类型
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenMalformed
			}
		}
		return nil, err
	}

	// 类型断言为自定义Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
