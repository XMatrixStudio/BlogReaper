package model

import (
	"github.com/boltdb/bolt"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
)

type CategoryModel struct {
	*Model
}

type Category struct {
	ID   bson.ObjectId `bson:"id"`   // 订阅分类的ID
	Name string        `bson:"name"` // 订阅分类的名称
}

func (m *CategoryModel) AddCategory(userID, name string) (category Category, err error) {
	return category, m.Update(func(b *bolt.Bucket) error {
		ub, err := b.CreateBucketIfNotExists([]byte(userID))
		if err != nil {
			return err
		}
		nub, err := ub.CreateBucketIfNotExists([]byte("key_name_value_userId"))
		if err != nil {
			return err
		}
		if nub.Get([]byte(name)) != nil {
			return errors.New("repeat_name")
		}
		category = Category{
			ID:   bson.NewObjectId(),
			Name: name,
		}
		bytes, err := bson.Marshal(&category)
		if err != nil {
			return err
		}
		err = ub.Put([]byte(category.ID), bytes)
		if err != nil {
			return err
		}
		return nub.Put([]byte(name), []byte(category.ID))
	})
}

func (m *CategoryModel) GetCategories(userID string) (categories []Category, err error) {
	return categories, m.View(func(b *bolt.Bucket) error {
		ub := b.Bucket([]byte(userID))
		if ub == nil {
			return nil
		}
		return ub.ForEach(func(k, v []byte) error {
			if string(k) != "key_name_value_userId" {
				category := Category{}
				err = bson.Unmarshal(v, &category)
				if err != nil {
					return err
				}
				categories = append(categories, category)
			}
			return nil
		})
	})
}

func (m *CategoryModel) EditCategory(userID, categoryID, newName string) (success bool, err error) {
	panic("not implement")
}
