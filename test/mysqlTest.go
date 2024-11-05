package test

import (
	"database/sql"
	"fmt"
	"goUserCenter/gxmodule"
	"log"
)

func M() {
	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建新用户示例
	newUser := gxmodule.NewUser("test1", "test1", "test1@example.com")
	if err := newUser.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	if err := newUser.Save(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("新用户创建成功，ID:", newUser.ID)

	// 用户登录验证示例
	inputUsername := "test1"
	inputPassword := "test1"
	user := gxmodule.User{Username: inputUsername, Password: inputPassword}
	if user.Authenticate(db) {
		fmt.Println("登录成功！")
	} else {
		fmt.Println("登录失败！")
	}

	// 更新用户信息示例
	userToUpdate := gxmodule.User{ID: newUser.ID, Username: "updated_user", Password: "updated_password", Email: "updated_user@example.com"}
	if err := userToUpdate.Update(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("用户信息更新成功！")

	// 切换用户账号激活状态示例
	if err := userToUpdate.ToggleActivity(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("用户账号激活状态切换成功！")
}

// 连接数据库并返回数据库连接对象
func connectDB() (*sql.DB, error) {
	// 根据实际情况修改这里的连接字符串
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gxUserCenter?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	// 尝试连接数据库，检查连接是否成功
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
