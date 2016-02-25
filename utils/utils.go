package utils

import (
	"time"
	crand "crypto/rand"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v8"
	"io"
	"log"
	"golang.org/x/crypto/scrypt"
	"fmt"
	"regexp"
	"math/rand"
)

var validate *validator.Validate

const (
// PwSaltBytes ...
	PwSaltBytes = 32
// PwHashBytes ...
	PwHashBytes = 64
)

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}
//SigninType 验证登录用户名类型
func SigninType(name string) string {
	if err := validate.Field(name, "required,email"); err == nil {
		return "email"
	}
	if b, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, name); b {
		return "mobile"
	}
	return ""
}

// JwtToken generate toke by signingKey
func JwtToken(signingKey string, claims map[string]interface{}) string {
	// New web token.
	token := jwt.New(jwt.SigningMethodHS256)

	// Set a header and a claim
	token.Header["typ"] = "JWT"   //加密方式
	token.Header["alg"] = "HS256" //加密类型
	if claims == nil {
		claims = map[string]interface{}{}
	}
	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(time.Hour * 96).Unix() //过期时间
	}
	token.Claims = claims
	// Generate encoded token
	t, _ := token.SignedString([]byte(signingKey))
	return t
}

// Salt 生成加密盐
func Salt() []byte {
	salt := make([]byte, PwSaltBytes)
	_, err := io.ReadFull(crand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	return salt
}

// Hash 密码加密
func Hash(salt, password []byte) []byte {
	hash, err := scrypt.Key([]byte(password), salt, 1 << 14, 8, 1, PwHashBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", hash)
	return hash
}

// RandNum return length 6
func RandNum() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100000)
}
