package main

import (
	"database/sql"
	"log"
)

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

func main() {
	//test.M()
	//test.G()
	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 设置路由
	router := SetupRoutes(db)

	// 启动应用
	router.Run(":8080")
}
