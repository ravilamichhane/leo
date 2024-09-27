package hasher

import (
	"github.com/ravilmc/leo/web"

	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHash(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}

func GenerateOTPHash(len int) (string, string, error) {
	otp := web.GenerateOtp(len)
	hash, err := Hash(otp)
	return otp, hash, err
}
