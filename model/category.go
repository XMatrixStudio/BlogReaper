package model

type CategoryModel struct {
	*Model
}

type Category struct {
	ID    string `bson:"id"`	// 订阅分类的ID
	Name  string `bson:"name"`	// 订阅分类的名称
	Feeds []Feed `bson:"feeds"`	// 分类中的订阅源
}
