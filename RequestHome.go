package main

import (
	"GoWeb/Error"
	"GoWeb/File"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func validateLogin(user *File.User, passwd string) error {

	if user == nil {
		return errors.New(Error.InvalidLogin_UserDoesNotExist)
	}

	if len(user.UserName) == 0 || len(passwd) == 0 {
		return errors.New(Error.InvalidLogin_WrongUserNameOrPassword)
	}

	hash := sha256.New()
	hash.Write([]byte(passwd))
	hashResult := hex.EncodeToString(hash.Sum(nil))
	hashResult = strings.ToUpper(hashResult)
	if user.PasswordHash != hashResult {
		return errors.New(Error.InvalidLogin_WrongPassword)
	}

	return nil
}

func validateMainPageInput(r *http.Request) (map[string]float64, error) {
	if strings.TrimSpace(r.Form["option1"][0]) == "" || strings.TrimSpace(r.Form["option2"][0]) == "" ||
		strings.TrimSpace(r.Form["option3"][0]) == "" {
		return nil, errors.New(Error.InvalidInput_RequiredOption)
	}

	if strings.TrimSpace(r.Form["option1_weight"][0]) == "" || strings.TrimSpace(r.Form["option2_weight"][0]) == "" ||
		strings.TrimSpace(r.Form["option3_weight"][0]) == "" {
		return nil, errors.New(Error.InvalidInput_RequiredOptionWeight)
	}

	weightMap, err := transformWeightToSum(r)
	if err != nil {
		return nil, err
	}

	weightSum := findSumOfWeight(weightMap)
	if weightSum <= 0 || weightSum > 100 {
		return nil, errors.New(Error.InvalidInput_InvalidSumOfWeight)
	}

	return weightMap, nil
}

func transformWeightToSum(r *http.Request) (map[string]float64, error) {
	resultArr := make(map[string]float64)

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

func transFormAndCheckInvalidNumber(weight string) (float64, error) {
	result, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return -1, err
	}
	if result <= 0 {
		return -1, errors.New(Error.InvalidInput_NegativeOrZeroWeight)
	}

	return result, nil
}

func isBothOptionAndWeightNotEmpty(option string, weight string) bool {
	return strings.TrimSpace(option) != "" && strings.TrimSpace(weight) != ""
}

func findSumOfWeight(weights map[string]float64) float64 {
	var sum float64 = 0
	for _, weight := range weights {
		sum += weight
	}
	return sum
}

func validateResultStatus(resultsNearBy []*File.NearBy) error {
	var allZeroResult = true
	for _, result := range resultsNearBy {
		if result.Status == File.STATUS_OK {
			allZeroResult = false
		}

		if result.Status != File.STATUS_OK && result.Status != File.STATUS_ZERO_RESULTS {
			return errors.New(Error.OutputError_RequestError)
		}
	}

	if allZeroResult {
		return errors.New(Error.OutputError_AllZeroResults)
	}

	return nil
}

func eliminateNotOpenResult(resultsNearBy []*File.NearBy) {
	for _, result := range resultsNearBy {
		stores := result.NearByResults
		resultStore := make([]File.NearByResult, 0)
		for i := 0; i < len(stores); i++ {
			if (stores[i].BusinessStatus == File.BUSINESS_STATUS_OPERATIONAL) && stores[i].OpeningHours.OpenNow == true {
				resultStore = append(resultStore, stores[i])
			} else {
				fmt.Println(stores[i])
			}
		}
		result.NearByResults = resultStore
	}
}

func rankAllResults(resultsNearBy []*File.NearBy, optionWeightMap *map[string]float64) []File.NearByResult {
	ret := make([]File.NearByResult, 0)

	for _, result := range resultsNearBy {
		for _, resultStore := range result.NearByResults {
			ret = append(ret, resultStore)
		}
	}

	return ret
}

func computeRankingScore(resultStore *File.NearByResult, weight float64) float64 {

}
