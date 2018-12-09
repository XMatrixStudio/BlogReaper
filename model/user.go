package model

import (
	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id"`   // 用户ID
	VioletID bson.ObjectId `bson:"vid"`   // VioletID
	Token    string        `bson:"token"` // Violet 访问令牌
	Email    string        `bson:"email"` // 用户唯一邮箱
	Class    int           `bson:"class"` // 用户类型
	Info     UserInfo      `bson:"info"`  // 用户个性信息
}

// UserInfo 用户个性信息
type UserInfo struct {
	Name   string `bson:"name"`   // 用户昵称
	Avatar string `bson:"avatar"` // 头像URL
	Bio    string `bson:"bio"`    // 个人简介
	Gender int    `bson:"gender"` // 性别
}
