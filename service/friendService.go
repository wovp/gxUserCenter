package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goUserCenter/gxmodule"
	"strconv"
	"time"
)

// SendFriendRequestHandler处理发送好友请求的函数
func SendFriendRequestHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token或其他方式获取当前用户ID（这里假设从token中获取，你可以根据实际情况调整）
		currentUserID := getUserIdFromToken(c)

		// 获取要发送好友请求的目标用户的用户名（假设通过表单提交）
		targetUsername := c.PostForm("target_username")

		// 根据用户名查询目标用户的ID
		targetUserID, err := gxmodule.GetUserIDByUsername(db, targetUsername)
		if err != nil {
			c.JSON(400, gin.H{"error": "找不到指定的目标用户"})
			return
		}

		// 创建好友关系实例，并设置请求状态为pending
		friendship := gxmodule.NewFriendship(currentUserID, targetUserID)
		if friendship.CheckRepeatedFriend(db) {
			c.JSON(500, gin.H{"message": "重复添加好友"})
			return
		}

		// 保存好友关系（包含请求状态等信息）到数据库
		if err := friendship.Save(db); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "好友请求已发送"})
	}
}

// AcceptFriendRequestHandler处理同意好友请求的函数
func AcceptFriendRequestHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token或其他方式获取当前用户ID（这里假设从token中获取，你可以根据实际情况调整）
		currentUserID := getUserIdFromToken(c)

		// 获取要同意的好友请求的ID（假设通过表单提交或其他方式获取）
		requestID := c.PostForm("request_id")

		// 根据请求ID查询对应的好友关系记录
		friendship, err := gxmodule.GetFriendshipByRequestID(db, strconv.Itoa(currentUserID), requestID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 检查当前用户是否是被请求的对象（即friend_id是否等于当前用户ID）
		if friendship.FriendID != currentUserID {
			c.JSON(400, gin.H{"error": "您无权同意此好友请求"})
			return
		}

		// 更新好友关系记录的请求状态为accepted，并更新请求接受时间
		friendship.RequestStatus = "accepted"
		friendship.RequestAcceptedAt = time.Now()

		// 更新好友关系记录到数据库
		if err := friendship.Update(db); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "好友请求已同意，双方已成为好友"})
	}
}

// GetFriendListHandler获取好友列表的函数
func GetFriendListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token或其他方式获取当前用户ID（这里假设从token中获取，你可以根据实际情况调整）
		currentUserID := getUserIdFromToken(c)

		// 获取当前用户的所有好友关系
		friendships, err := gxmodule.GetFriendshipsByUserID(db, currentUserID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 根据好友关系中的用户ID获取好友详细信息（这里假设已经有一个函数GetUserDetailByUserID用于获取用户详细信息）
		var friendDetails []*gxmodule.UserDetail
		for _, friendship := range friendships {
			var friendID int
			if friendship.UserID == currentUserID {
				friendID = friendship.FriendID
			} else {
				friendID = friendship.UserID
			}

			friendDetail, err := gxmodule.GetUserDetailByUserID(db, friendID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			friendDetails = append(friendDetails, friendDetail)
		}

		c.JSON(200, gin.H{"friend_list": friendDetails})
	}
}

// GetPendingFriendRequestsHandler处理查询待处理好友请求列表的函数
func GetPendingFriendRequestsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token或其他方式获取当前用户ID（这里假设从token中获取，你可以根据实际情况调整）
		currentUserID := getUserIdFromToken(c)

		// 获取当前用户的所有待处理好友请求
		pendingFriendRequests, err := gxmodule.GetFriendshipsPendingByUserID(db, currentUserID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 根据好友请求中的用户ID获取好友详细信息（这里假设已经有一个函数GetUserDetailByUserID用于获取用户详细信息）
		var friendDetails []*gxmodule.UserDetail
		for _, request := range pendingFriendRequests {
			var friendID int
			if request.UserID == currentUserID {
				friendID = request.FriendID
			} else {
				friendID = request.UserID
			}

			friendDetail, err := gxmodule.GetUserDetailByUserID(db, friendID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			friendDetails = append(friendDetails, friendDetail)
		}

		c.JSON(200, gin.H{"pending_friend_requests": friendDetails})
	}
}
