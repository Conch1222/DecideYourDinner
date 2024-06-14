package main

import (
	"GoWeb/Error"
	"GoWeb/File"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func validateLogin(user *File.User, passwd string) (bool, error) {

	if user == nil {
		return false, errors.New(Error.InvalidLogin_UserDoesNotExist)
	}

	if len(user.UserName) == 0 || len(passwd) == 0 {
		return false, errors.New(Error.InvalidLogin_WrongUserNameOrPassword)
	}

	hash := sha256.New()
	hash.Write([]byte(passwd))
	hashResult := hex.EncodeToString(hash.Sum(nil))
	hashResult = strings.ToUpper(hashResult)
	if user.PasswordHash != hashResult {
		return false, errors.New(Error.InvalidLogin_WrongPassword)
	}

	return true, nil
}

func validateMainPageInput(r *http.Request) (bool, error) {
	if r.Form["option1"][0] == "" || r.Form["option2"][0] == "" || r.Form["option3"][0] == "" {
		return false, errors.New(Error.InvalidInput_RequiredOption)
	}

	if r.Form["option1_wight"] == nil || r.Form["option2_wight"] == nil || r.Form["option3_wight"] == nil {
		return false, errors.New(Error.InvalidInput_RequiredOptionWeight)
	}

	return true, nil
}
