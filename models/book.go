package models

type Book struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	AuthorName  string  `json:"author_name"`
	Ebook       bool    `json:"ebook"`
	PublishDate string  `json:"publish_date"`
}

type ShortBook struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type PriceRange struct {
	low  float64
	high float64
}

type BookQuery struct {
	Title      string `json:"title"`
	PriceRange string `json:"price_range"`
	AuthorName string `json:"author_name"`
}

type Store struct {
	Books   int64 `json:"books"`
	Authors int64 `json:"authors"`
}
