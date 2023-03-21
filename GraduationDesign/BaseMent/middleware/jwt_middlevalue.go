package middleware

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JWT struct {
	JwtKey []byte
}

func NewJwt() *JWT {
	return &JWT{[]byte(Config.Conf.JwtK.JwtKey)}
}

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

// 定义错误
var (
	TokenExpired     error = errors.New("token已过期,请重新登录")
	TokenNotValueYet error = errors.New("此token无效,请重新登录")
	TokenMalFormed   error = errors.New("token不正确,请重新登录")
	TokenInvalid     error = errors.New("token格式错误,请重新登录")
)

// CreatToken CreatToken
func (j *JWT) CreatToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParserToken ParserToken
func (j *JWT) ParserToken(tokenString string) (*Claims, error) {
	var JwtKey = []byte(Config.Conf.JwtK.JwtKey)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalFormed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValueYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claimsData, ok := token.Claims.(*Claims); ok && token.Valid {
			return claimsData, nil
		}
	}
	return nil, TokenInvalid
}

// JwtToken middleValue
func JwtToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		var code int
		//store in Header of request
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			code = utils.ErrorTokenExist
			context.JSON(http.StatusUnauthorized, gin.H{
				"status":  code,
				"message": utils.GetErrMsg(code),
			})
			context.Abort()
			return
		}
		checkToken := strings.Split(tokenString, " ")
		if len(checkToken) == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"status":  code,
				"message": utils.GetErrMsg(code),
			})
			context.Abort()
			return
		}
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"status":  code,
				"message": utils.GetErrMsg(code),
			})
			context.Abort()
			return
		}

		j := NewJwt()

		//解析Token
		claims, err := j.ParserToken(checkToken[1])
		if err != nil {
			if err == TokenExpired {
				context.JSON(http.StatusUnauthorized, gin.H{
					"status":  utils.ERROR,
					"message": TokenExpired,
					"data":    nil,
				})
				context.Abort()
				return
			}
			context.JSON(http.StatusUnauthorized, gin.H{
				"status":  utils.ERROR,
				"message": err.Error(),
				"data":    TokenExpired,
			})
			context.Abort()
			return
		}
		context.Set("user_id", claims.UserId)
		context.Set("claims", claims)
		context.Next()
	}
}
