package gxmodule

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"goUserCenter/config"
	"time"
)

// User User结构体表示用户信息
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

// NewUser NewUser创建一个新的User结构体实例
func NewUser(username, password, email string) *User {
	return &User{
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}
}

// Save Save将用户信息保存到数据库中
func (u *User) Save(db *sql.DB) error {
	if u.ID == 0 {
		// 插入新用户
		query := "INSERT INTO users (username, password, email, created_at, updated_at, is_active) VALUES (?,?,?,?,?,?)"
		result, err := db.Exec(query, u.Username, u.Password, u.Email, u.CreatedAt, u.UpdatedAt, u.IsActive)
		if err != nil {
			return err
		}

		// 获取插入后新用户的ID
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		u.ID = int(id)
	} else {
		// 更新现有用户
		query := "UPDATE users SET username =?, password =?, email =?, updated_at =? WHERE id =?"
		_, err := db.Exec(query, u.Username, u.Password, u.Email, time.Now(), u.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate Validate验证用户输入信息的有效性
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if u.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if u.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// 这里可以添加更多复杂的验证逻辑，比如邮箱格式是否正确等

	return nil
}

// Update Update更新用户信息
func (u *User) Update(db *sql.DB) error {
	query := "UPDATE users SET username =?, password =?, email =?, updated_at =? WHERE id =?"
	_, err := db.Exec(query, u.Username, u.Password, u.Email, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ToggleActivity ToggleActivity切换用户账号的激活状态
func (u *User) ToggleActivity(db *sql.DB) error {
	u.IsActive = !u.IsActive
	query := "UPDATE users SET is_active =? WHERE id =?"
	_, err := db.Exec(query, u.IsActive, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate验证用户登录信息并生成 JWT token
func (u *User) Authenticate(db *sql.DB) (string, bool) {
	query := "SELECT id, username, password, email, created_at, updated_at, is_active FROM users WHERE username =? AND password =?"
	row := db.QueryRow(query, u.Username, u.Password)

	var storedUser User
	err := row.Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password, &storedUser.Email, &storedUser.CreatedAt, &storedUser.UpdatedAt, &storedUser.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		return "", false
	}

	// 如果登录成功，生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(config.AppConfig.TokenExpiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.SecretKey))
	if err != nil {
		return "", false
	}

	return tokenString, true
}

// CreateTable 创建用户表
func (u *User) CreateTable(db *sql.DB) error {
	createUserTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL,
        is_active BOOLEAN NOT NULL
    );
    `
	_, err := db.Exec(createUserTableSQL)
	return err
}
