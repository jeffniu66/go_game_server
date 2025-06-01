package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// 三个参数，分别是原始密码、盐、迭代次数
func PasswordEncode(password string, salt string, iterations int32) (string, error) {
	// 如果没有设置盐，则使用12位的随机字符串
	if strings.TrimSpace(salt) == "" {
		salt = RandStr(6)
	}

	// 确保盐不包含美元$符号
	if strings.Contains(salt, "$") {
		return "", errors.New("salt contains dollar sign ($)")
	}

	// 如果迭代次数小于等于0，则设置为8000
	if iterations <= 0 {
		iterations = 8000
	}

	// pbkdf2加密
	hash := pbkdf2.Key([]byte(password), []byte(salt), int(iterations), sha256.Size, sha256.New)

	b64Hash := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", iterations, salt, b64Hash), nil
}

func PasswordVerify(password string, encoded string) (bool, error) {
	// 先根据美元$符号分割密钥为4个子字符串
	s := strings.Split(encoded, "$")

	// 如果分割结果不是4个子字符串，则认为不是pbkdf2_sha256算法的结果密钥，跳出错误
	if len(s) != 4 {
		return false, errors.New("hashed password components mismatch")
	}

	// 分割子字符串的结果分别为算法名、迭代次数、盐和base64编码
	algorithm, iterations, salt := s[0], s[1], s[2]

	// 如果密钥算法名不是pbkdf2_sha256算法，跳出错误
	if algorithm != "pbkdf2_sha256" {
		return false, errors.New("algorithm mismatch")
	}

	// 加密用的迭代次数
	i := ToInt(iterations)

	// 将原始密码用上面获取的盐、迭代次数进行加密
	newEncoded, err := PasswordEncode(password, salt, i)
	if err != nil {
		return false, err
	}

	// 最终用hmac.Equal函数判断两个密钥是否相同
	return hmac.Equal([]byte(newEncoded), []byte(encoded)), nil
}
