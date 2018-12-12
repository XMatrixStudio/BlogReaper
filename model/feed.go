package model

type FeedModel struct {
	*Model
}

type Feed struct {
	URL        string    `bson:"url"`			// 订阅源的URL
	Title      string    `bson:"title"`			// 订阅源的标题
	Subtitle   string    `bson:"subtitle"`		// 订阅源的子标题
	Articles   []Article `bson:"articles"`		// 订阅源包括的文章
}
