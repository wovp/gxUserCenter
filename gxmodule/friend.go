package gxmodule

import (
	"database/sql"
	"fmt"
	"time"
)

// Friendship结构体表示好友关系
type Friendship struct {
	ID                int
	UserID            int
	FriendID          int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RequestStatus     string
	RequestSentAt     time.Time
	RequestAcceptedAt time.Time
}

// NewFriendship创建一个新的Friendship结构体实例
func NewFriendship(userID, friendID int) *Friendship {
	return &Friendship{
		UserID:            userID,
		FriendID:          friendID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		RequestStatus:     "pending",
		RequestSentAt:     time.Now(),
		RequestAcceptedAt: time.Now(),
	}
}

// Save将好友关系保存到数据库中
func (f *Friendship) Save(db *sql.DB) error {
	query := `
    INSERT INTO friendships (user_id, friend_id, created_at, updated_at, request_status, request_sent_at, request_accepted_at)
    VALUES (?,?,?,?,?,?,?)
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(f.UserID, f.FriendID, f.CreatedAt, f.UpdatedAt, f.RequestStatus, f.RequestSentAt, f.RequestAcceptedAt)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

// GetFriendshipsByUserID根据用户ID获取该用户的所有好友关系
func GetFriendshipsByUserID(db *sql.DB, userID int) ([]*Friendship, error) {
	query := `
    SELECT id, user_id, friend_id, created_at, updated_at, request_status, request_sent_at, request_accepted_at
    FROM friendships
    WHERE (user_id =? AND request_status = 'accepted')  OR (friend_id =? AND request_status = 'accepted') 
    `
	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*Friendship
	for rows.Next() {
		var friendship Friendship
		err := rows.Scan(&friendship.ID, &friendship.UserID, &friendship.FriendID, &friendship.CreatedAt, &friendship.UpdatedAt, &friendship.RequestStatus, &friendship.RequestSentAt, &friendship.RequestAcceptedAt)
		if err != nil {
			return nil, err
		}

		friendships = append(friendships, &friendship)
	}

	return friendships, nil
}

// GetFriendshipsPendingByUserID根据用户ID获取该用户的所有待处理的好友关系请求(其他人添加该用户)
func GetFriendshipsPendingByUserID(db *sql.DB, userID int) ([]*Friendship, error) {
	query := `
    SELECT id, user_id, friend_id, created_at, updated_at, request_status, request_sent_at, request_accepted_at
    FROM friendships
    WHERE (friend_id =? AND request_status = 'pending')
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*Friendship
	for rows.Next() {
		var friendship Friendship
		err := rows.Scan(&friendship.ID, &friendship.UserID, &friendship.FriendID, &friendship.CreatedAt, &friendship.UpdatedAt, &friendship.RequestStatus, &friendship.RequestSentAt, &friendship.RequestAcceptedAt)
		if err != nil {
			return nil, err
		}

		friendships = append(friendships, &friendship)
	}

	return friendships, nil
}

// GetFriendshipByRequestID根据请求ID查询对应的好友关系记录
func GetFriendshipByRequestID(db *sql.DB, currentUserID string, requestID string) (*Friendship, error) {
	query := `
    SELECT id, user_id, friend_id, created_at, updated_at, request_status, request_sent_at, request_accepted_at
    FROM friendships
    WHERE user_id =? AND friend_id =?
    `
	// re 请求 cur
	row := db.QueryRow(query, requestID, currentUserID)

	var friendship Friendship
	err := row.Scan(&friendship.ID, &friendship.UserID, &friendship.FriendID, &friendship.CreatedAt, &friendship.UpdatedAt, &friendship.RequestStatus, &friendship.RequestSentAt, &friendship.RequestAcceptedAt)
	if err != nil {
		return nil, err
	}

	return &friendship, nil
}

// Update更新Friendship结构体实例到数据库中的记录
func (f *Friendship) Update(db *sql.DB) error {
	query := `
    UPDATE friendships
    SET user_id =?, friend_id =?, created_at =?, updated_at =?, request_status =?, request_sent_at =?, request_accepted_at =?
    WHERE id =?
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(f.UserID, f.FriendID, f.CreatedAt, f.UpdatedAt, f.RequestStatus, f.RequestSentAt, f.RequestAcceptedAt, f.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("没有找到对应的好友关系记录进行更新")
	}

	return nil
}

func (f *Friendship) CheckRepeatedFriend(db *sql.DB) bool {
	query := `
    SELECT id, user_id, friend_id, created_at, updated_at, request_status, request_sent_at, request_accepted_at
    FROM friendships
    WHERE (user_id = ? AND friend_id = ?) OR  (user_id = ? AND friend_id = ?)
    `

	rows, err := db.Query(query, f.UserID, f.FriendID, f.FriendID, f.UserID)

	if err != nil {
		return false
	}
	defer rows.Close()
	// 如果查询结果有行数据，说明存在重复的好友关系
	return rows.Next()

}
