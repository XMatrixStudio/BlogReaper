package model

type ArticleModel struct {
	*Model
}

type Article struct {
	URL        string   `bson:"url"`
	Title      string   `bson:"title"`
	Published  string   `bson:"published"`
	Updated    string   `bson:"updated"`
	Content    string   `bson:"content"`
	Summary    string   `bson:"summary"`
	Categories []string `bson:"categories"`
	Read       bool     `bson:"read"`
	Later      bool     `bson:"later"`
}
