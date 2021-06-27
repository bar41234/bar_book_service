package models

type Book struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	AuthorName  string  `json:"author_name"`
	Ebook       bool    `json:"ebook"`
	PublishDate string  `json:"publish_date"`
}

type Store struct {
	Books   int `json:"books"`
	Authors int `json:"authors"`
}
