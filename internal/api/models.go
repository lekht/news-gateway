package api

type NewFullDetailed struct {
	ID       int
	Title    string
	PubTime  int64
	Link     string
	Content  string
	Comments []Comment
}

type NewsShortDetailed struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}

type NewsList struct {
	Posts []NewsShortDetailed
	Page  Page
}

type Comment struct {
	ID       int    `json:"id"`
	NewsID   int    `json:"news_id"`
	ParentID int    `json:"parent_id"`
	Msg      string `json:"msg"`
	PubTime  int64  `json:"pub_time"`
}

type Page struct {
	TotalPages  int
	CurrentPage int
	NewsPerPage int
}
