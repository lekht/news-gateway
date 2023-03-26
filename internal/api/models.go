package api

type NewFullDetailed struct {
	ID       int       // номер записи
	Title    string    // заголовок публикации
	PubTime  int64     // время публикации
	Link     string    // ссылка на источник
	Content  string    // содержание публикации
	Comments []Comment // комментарии к публикации
}

type NewsShortDetailed struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

type NewsList struct {
	Posts []NewsShortDetailed // список сокращенных новостей
	Page  Page
}

type Comment struct {
	ID       int
	NewsID   int // id новости
	ParentID int // id родительского комментария
	Text     string
	PubTime  int64
}

type Page struct {
	TotalPages  int
	CurrentPage int
	NewsPerPage int
}
