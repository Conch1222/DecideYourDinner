package main

import (
	"GoWeb/Type"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func (dbConn *DBConnection) getUserByUserName(userName string) *Type.User {
	user := new(Type.User)
	row := dbConn.db.QueryRow("select user_id, user_name, password_hash from web_user where user_name = ?", userName)
	if err := row.Scan(&user.ID, &user.UserName, &user.PasswordHash); err != nil {
		fmt.Printf("mapping error: %v\n", err)
		return nil
	}
	fmt.Println("success: ", *user)
	return user
}

func (dbConn *DBConnection) isUserNameDuplicated(userName string) bool {
	count := 0
	row := dbConn.db.QueryRow("select count(user_id) from web_user where user_name = ?", userName)
	if err := row.Scan(&count); err != nil {
		return true
	}
	return count > 0
}

func (dbConn *DBConnection) saveUserInfo(userName string, password string) (int64, error) {

	hash := sha256.New()
	hash.Write([]byte(password))
	hashResult := hex.EncodeToString(hash.Sum(nil))
	hashResult = strings.ToUpper(hashResult)

	result, err := dbConn.db.Exec("insert into web_user (user_name, password_hash, create_time) values (?, ?, ?)", userName, hashResult, time.Now())
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (dbConn *DBConnection) saveQueryRecord(userId int64, storeName string, storeAddress string, storeRating float64, storeMapLink string) (int64, error) {
	query := "insert into web_query_record (user_id, store_name, store_address, store_rating, store_map_link, create_time) values (?, ?, ?, ?, ?, ?)"

	result, err := dbConn.db.Exec(query, userId, storeName, storeAddress, storeRating, storeMapLink, time.Now())
	if err != nil {
		return 0, err
	}

	queryRecordId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return queryRecordId, nil
}

func (dbConn *DBConnection) getQueryRecord(userId int64, recordNumber int) ([]Type.QueryData, error) {
	query := "select store_name, store_address, store_rating, store_map_link from web_query_record where user_id = ? order by create_time desc limit ?;"

	rows, err := dbConn.db.Query(query, userId, recordNumber)
	if err != nil {
		return nil, err
	}

	var results []Type.QueryData

	for rows.Next() {
		var result Type.QueryData
		err := rows.Scan(&result.StoreName, &result.StoreAddress, &result.StoreRating, &result.StoreMapLink)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
