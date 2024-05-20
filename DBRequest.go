package main

import (
	"database/sql"
	"fmt"
)

func getUserByUserName(db *sql.DB, userName string) *User {
	user := new(User)
	row := db.QueryRow("select * from web_user where user_name = ?", userName)
	if err := row.Scan(&user.ID, &user.lastName, &user.firstName, &user.userName, &user.passwordHash, &user.createTime); err != nil {
		fmt.Printf("mapping error: %v\n", err)
		return nil
	}
	fmt.Println("success: ", *user)
	return user
}
