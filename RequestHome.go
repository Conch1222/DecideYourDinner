package main

import (
	"GoWeb/Error"
	"GoWeb/Type"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func validateSignUp(userName string, password string, confirmPassword string) error {
	if len(strings.Trim(userName, " ")) == 0 {
		return errors.New(Error.InvalidSignUp_UserNameCannotBeEmpty)
	}

	if strings.Contains(userName, " ") {
		return errors.New(Error.InvalidSignUp_UserNameCannotContainSpaces)
	}

	if len(password) < 8 {
		return errors.New(Error.InvalidSignUp_PasswordTooShort)
	}

	if password != confirmPassword {
		return errors.New(Error.InvalidSignUp_PasswordDiffFromConfirmPassword)
	}

	DBConn := connectDB()
	if DBConn.isUserNameDuplicated(userName) {
		return errors.New(Error.InvalidSignUp_UserNameAlreadyTaken)
	}

	return nil
}

func validateLogin(user *Type.User, passwd string) error {

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

func validateResultStatus(resultsNearBy []*Type.NearBy) error {
	var allZeroResult = true
	for _, result := range resultsNearBy {
		if result.Status == Type.STATUS_OK {
			allZeroResult = false
		}

		if result.Status != Type.STATUS_OK && result.Status != Type.STATUS_ZERO_RESULTS {
			return errors.New(Error.OutputError_RequestError)
		}
	}

	if allZeroResult {
		return errors.New(Error.OutputError_AllZeroResults)
	}

	return nil
}

func eliminateNotOpenResult(resultsNearBy []*Type.NearBy) {
	for _, result := range resultsNearBy {
		stores := result.NearByResults
		resultStore := make([]Type.NearByResult, 0)
		for i := 0; i < len(stores); i++ {
			if (stores[i].BusinessStatus == Type.BUSINESS_STATUS_OPERATIONAL) && stores[i].OpeningHours.OpenNow == true {
				resultStore = append(resultStore, stores[i])
			} else {
				fmt.Println("eliminate not open")
			}
		}
		result.NearByResults = resultStore
	}
}

func rankAllResults(resultsNearBy []*Type.NearBy) []Type.NearByResult {
	ret := make([]Type.NearByResult, 0)

	for _, result := range resultsNearBy {
		for _, resultStore := range result.NearByResults {
			resultStore.RankingScore = computeRankingScore(&resultStore, result.Weight)
			ret = append(ret, resultStore)
		}
	}

	sort.Slice(ret, func(i, j int) bool { return ret[i].RankingScore > ret[j].RankingScore })
	return ret
}

func computeRankingScore(resultStore *Type.NearByResult, weight float64) float64 {
	priceScore := Type.Default_Price
	if resultStore.PriceLevel != 0 {
		priceScore = float64(4 - resultStore.PriceLevel)
	}

	ratingScore := Type.Default_Rating
	if resultStore.Rating > 1e-9 {
		ratingScore = resultStore.Rating
	}

	// formula: log(ratingScore) / 5.0
	userRatingScore := Type.Default_UserRating
	if resultStore.UserRatingTotal != 0 {
		userRatingScore = float64(resultStore.UserRatingTotal)
	}
	userRatingScore = math.Log(userRatingScore+1) / 5.0

	normalizedWeight := weight / 100.0

	return (Type.Weight_Price * priceScore) + (Type.Weight_Rating * ratingScore) +
		(Type.Weight_UserRating * ratingScore * userRatingScore) + (Type.Weight_UserWeight * normalizedWeight)
}

func convertAddress(store Type.NearByResult) string {
	var sb strings.Builder
	area := store.PlusCode.CompoundCode
	compoundCodeSplit := strings.Split(area, " ")
	area = compoundCodeSplit[len(compoundCodeSplit)-1]

	sb.WriteString(area)
	sb.WriteString(store.Vicinity)

	return sb.String()
}

func ReadKey(filePath string, error string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", errors.New(error)
}
