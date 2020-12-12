package utils

import "golang.org/x/crypto/bcrypt"

var PwdUtils pwdUtilsInterface = &pwdUtilsStruct{}

type pwdUtilsStruct struct{}

type pwdUtilsInterface interface {
	GeneratePwdHash(password string) (string, error)
}

func (pwdUtilsStruct) GeneratePwdHash(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}
