package echo

type UserInfoInRedis struct {
	Name              string `tomap:"name"`
	TotalArticleCount int    `tomap:"upload"`
	ID                uint   `tomap:"_id"`
	//Articles          []byte `tomap:"articles"`
	//Flag              int    `tomap:"flag"`
}
type ArticleInfoInRedis struct {
	ID      uint   `tomap:"id"`
	Title   string `tomap:"title"`
	Tags    string `tomap:"tags"`
	Content string `tomap:"content"`
}
