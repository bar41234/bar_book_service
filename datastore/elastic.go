package datastore

import (
	"context"
	"encoding/json"
	"github.com/bar41234/bar_book_service/models"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	indexName = "bar_books_store_2"
	querySize = 100
)

type bookManager struct {
	client *elastic.Client
}

type countDistinctAuthors struct {
	Val int `json:"value"`
}

func NewElastic(url string) (*bookManager, error) {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, errors.New("Elastic client creation was failed")
	}
	return &bookManager{client: client}, nil
}

func (e *bookManager) Get(id string) (*models.Book, error) {
	ctx := context.Background()
	res, err := e.client.Get().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	var book models.Book
	err = json.Unmarshal(res.Source, &book)
	if err != nil {
		return nil, err
	}
	book.Id = res.Id
	return &book, nil
}

func (e *bookManager) Add(book models.Book) (string, error) {
	ctx := context.Background()
	jsonBook, err := json.Marshal(book)
	if err != nil {
		return "", err
	}
	bookJason := string(jsonBook)
	res, err := e.client.Index().
		Index(indexName).
		BodyJson(bookJason).
		Do(ctx)
	if err != nil {
		return "", err
	}
	return res.Id, nil

}

func (e *bookManager) Update(id string, title string) (string, error) {
	ctx := context.Background()
	_, err := e.client.Update().Index(indexName).Id(id).Doc(map[string]interface{}{"title": title}).Do(ctx)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *bookManager) Delete(id string) error {
	ctx := context.Background()
	_, err := e.client.Delete().Index(indexName).Id(id).Do(ctx)
	return err
}

func (e *bookManager) searchQueryBuilder(title string, author string, priceRange string) (*elastic.BoolQuery, error) {
	query := elastic.NewBoolQuery()
	if title != "" {
		titleQuery := elastic.NewMatchQuery("title", title)
		query.Must(titleQuery)
	}
	if author != "" {
		authorQuery := elastic.NewMatchQuery("author_name", author)
		query.Must(authorQuery)
	}
	if priceRange != "" {
		rangeList := strings.Split(priceRange, ",")
		if len(rangeList) != 2 {
			err := errors.New("Price range argument is invalid! Please use the format 'lowRange,highRange'")
			return nil, err
		}
		low, err := strconv.ParseFloat(rangeList[0], 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(rangeList[1], 64)
		if err != nil {
			return nil, err
		}
		rangeQuery := elastic.NewRangeQuery("price").Gte(low).Lte(high)
		query.Must(rangeQuery)
	}
	return query, nil
}

func (e *bookManager) Search(title string, author string, priceRange string) ([]models.Book, error) {
	ctx := context.Background()
	query, err := e.searchQueryBuilder(title, author, priceRange)
	if err != nil {
		return nil, err
	}
	searchResult, err := e.client.Search().Index(indexName).Query(query).Size(querySize).Do(ctx)
	if err != nil {
		return nil, err
	}
	var allBooks []models.Book
	if searchResult.Hits.TotalHits.Value > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var book models.Book
			err := json.Unmarshal(hit.Source, &book)
			if err != nil {
				return nil, err
			}
			book.Id = hit.Id
			allBooks = append(allBooks, book)
		}
	}
	return allBooks, nil

}

func (e *bookManager) GetStore() (*models.Store, error) {
	ctx := context.Background()
	distinctAuthorsAgg := elastic.NewCardinalityAggregation().Field("author_name.keyword")
	allBooksSearchResults, err := e.client.Search().Index(indexName).Aggregation("distinct_authors", distinctAuthorsAgg).Size(0).Do(ctx)
	if err != nil {
		return nil, err
	}
	authorsCount, found := allBooksSearchResults.Aggregations["distinct_authors"]
	if !found {
		err = errors.New("Aggregation name was not found")
		return nil, &elastic.Error{Status: http.StatusNotFound}
	}
	count := countDistinctAuthors{}
	err = json.Unmarshal(authorsCount, &count)
	if err != nil {
		return nil, err
	}
	return &models.Store{
		Books:   int(allBooksSearchResults.Hits.TotalHits.Value),
		Authors: count.Val,
	}, nil
}
