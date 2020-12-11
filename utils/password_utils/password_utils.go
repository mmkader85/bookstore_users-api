package password_utils

import "golang.org/x/crypto/bcrypt"

func GeneratePwdHash(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}
