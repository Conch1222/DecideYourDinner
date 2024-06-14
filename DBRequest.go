package main

import (
	"GoWeb/File"
	"database/sql"
	"fmt"
)

func getUserByUserName(db *sql.DB, userName string) *File.User {
	user := new(File.User)
	row := db.QueryRow("select * from web_user where user_name = ?", userName)
	if err := row.Scan(&user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.PasswordHash, &user.CreateTime); err != nil {
		fmt.Printf("mapping error: %v\n", err)
		return nil
	}
	fmt.Println("success: ", *user)
	return user
}
