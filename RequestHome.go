package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

func validateLogin(user *User, passwd string) (bool, error) {

	if user == nil {
		return false, errors.New("user does not exist")
	}

	if len(user.userName) == 0 || len(passwd) == 0 {
		return false, errors.New("username or password is empty")
	}

	hash := sha256.New()
	hash.Write([]byte(passwd))
	hashResult := hex.EncodeToString(hash.Sum(nil))
	hashResult = strings.ToUpper(hashResult)
	if user.passwordHash != hashResult {
		return false, errors.New("password is wrong")
	}

	return true, nil
}
