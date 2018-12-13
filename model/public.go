package model

type PublicModel struct {
	*Model
}

type PublicFeed struct {
	URL        string          `bson:"url"`         // 订阅源的URL
	Title      string          `bson:"title"`       // 订阅源的标题
	Subtitle   string          `bson:"subtitle"`    // 订阅源的子标题
	Articles   []PublicArticle `bson:"articles"`    // 订阅源包括的文章
	Star       int             `bson:"star"`        // 订阅数量
	UpdateDate int64           `bson:"update_date"` // 更新时间，如果超过12小时就更新，或强制更新
}

type PublicArticle struct {
	URL        string   `bson:"url"`
	Title      string   `bson:"title"`
	Published  string   `bson:"published"`
	Updated    string   `bson:"updated"`
	Content    string   `bson:"content"`
	Summary    string   `bson:"summary"`
	Categories []string `bson:"categories"`
}
