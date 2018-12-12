package model

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/globalsign/mgo/bson"
)

type UserModel struct {
	*Model
}

type User struct {
	VioletID bson.ObjectId `bson:"vid"`   // VioletID
	Token    string        `bson:"token"` // Violet 访问令牌
	Email    string        `bson:"email"` // 用户唯一邮箱
	Class    int           `bson:"class"` // 用户类型
	Info     UserInfo      `bson:"info"`  // 用户个性信息
}

// UserInfo 用户个性信息
type UserInfo struct {
	Name   		string 		`bson:"name"`		// 用户昵称
	Avatar 		string 		`bson:"avatar"`		// 头像URL
	Bio    		string 		`bson:"bio"`		// 个人简介
	Gender 		int    		`bson:"gender"`		// 性别
	Categories 	[]Category 	`bson:"categories"`	// 订阅分类
}

// GetUserByID 根据ID查询用户
func (m *UserModel) GetUserByID(id string) (user User, err error) {
	if !bson.IsObjectIdHex(id) {
		return user, errors.New("not_id")
	}
	m.View(func(b *bolt.Bucket) error {
		bytes := b.Get([]byte(id))
		if bytes == nil {
			err = errors.New("not_found")
		} else {
			bson.Unmarshal(bytes, &user)
		}
		return nil
	})
	return
}

// AddUser 添加用户
func (m *UserModel) AddUser(id, token, email, name, avatar, bio string, gender int) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.Update(func(b *bolt.Bucket) error {
		bytes, err := bson.Marshal(User{
			VioletID: bson.ObjectIdHex(id),
			Token:    token,
			Email:    email,
			Class:    0,
			Info: UserInfo{
				Name:   name,
				Avatar: avatar,
				Bio:    bio,
				Gender: gender,
			},
		})
		if err != nil {
			return err
		}
		b.Put([]byte(id), bytes)
		return nil
	})
}

// SetUserToken 设置Token
func (m *UserModel) SetUserToken(id, token string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.Update(func(b *bolt.Bucket) error {
		bytes := b.Get([]byte(id))
		user := User{}
		bson.Unmarshal(bytes, &user)
		user.Token = token
		bytes, err := bson.Marshal(user)
		if err != nil {
			return err
		}
		b.Put([]byte(id), bytes)
		return nil
	})
}
