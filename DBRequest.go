package main

import (
	"GoWeb/File"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func getUserByUserName(db *sql.DB, userName string) *File.User {
	user := new(File.User)
	row := db.QueryRow("select user_name, password_hash from web_user where user_name = ?", userName)
	if err := row.Scan(&user.UserName, &user.PasswordHash); err != nil {
		fmt.Printf("mapping error: %v\n", err)
		return nil
	}
	fmt.Println("success: ", *user)
	return user
}

func isUserNameDuplicated(db *sql.DB, userName string) bool {
	count := 0
	row := db.QueryRow("select count(user_id) from web_user where user_name = ?", userName)
	if err := row.Scan(&count); err != nil {
		return true
	}
	return count > 0
}

func saveUserInfo(db *sql.DB, userName string, password string) (int64, error) {

	hash := sha256.New()
	hash.Write([]byte(password))
	hashResult := hex.EncodeToString(hash.Sum(nil))
	hashResult = strings.ToUpper(hashResult)

	result, err := db.Exec("insert into web_user (user_name, password_hash, create_time) values (?, ?, ?)", userName, hashResult, time.Now())
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}
