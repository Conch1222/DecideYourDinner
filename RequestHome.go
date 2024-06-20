package main

import (
	"GoWeb/Error"
	"GoWeb/File"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
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

func validateMainPageInput(r *http.Request) (bool, error, map[string]int) {
	if r.Form["option1"][0] == "" || r.Form["option2"][0] == "" || r.Form["option3"][0] == "" {
		return false, errors.New(Error.InvalidInput_RequiredOption), nil
	}

	if r.Form["option1_weight"][0] == "" || r.Form["option2_weight"][0] == "" || r.Form["option3_weight"][0] == "" {
		return false, errors.New(Error.InvalidInput_RequiredOptionWeight), nil
	}

	weightMap, err := transformWeightToSum(r)
	if err != nil {
		return false, err, nil
	}

	weightSum := findSumOfWeight(weightMap)
	if weightSum <= 0 || weightSum > 100 {
		return false, errors.New(Error.InvalidInput_InvalidSumOfWeight), nil
	}

	return true, nil, weightMap
}

func transformWeightToSum(r *http.Request) (map[string]int, error) {
	resultArr := make(map[string]int)

	for i := 1; i <= 5; i++ {
		optionNum := "option" + strconv.Itoa(i)
		weight := optionNum + "_weight"

		if i > 3 && !isBothOptionAndWeightNotEmpty(r.Form[optionNum][0], r.Form[weight][0]) {
			continue
		}

		result, err := transFormAndCheckInvalidNumber(r.Form[weight][0])
		if err != nil {
			return nil, err
		}

		if _, hasV := resultArr[r.Form[optionNum][0]]; hasV {
			return nil, errors.New(Error.InvalidInput_DuplicateOption)
		} else {
			resultArr[r.Form[optionNum][0]] = result
		}
	}
	return resultArr, nil
}

func transFormAndCheckInvalidNumber(weight string) (int, error) {
	result, err := strconv.Atoi(weight)
	if err != nil {
		return -1, err
	}
	if result <= 0 {
		return -1, errors.New(Error.InvalidInput_NegativeOrZeroWeight)
	}

	return result, nil
}

func isBothOptionAndWeightNotEmpty(option string, weight string) bool {
	return option != "" && weight != ""
}

func findSumOfWeight(weights map[string]int) int {
	sum := 0
	for _, weight := range weights {
		sum += weight
	}
	return sum
}
