package datastore

import (
	"context"
	"encoding/json"
	"github.com/bar41234/bar_book_service/models"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const (
	ELASTIC_URL = "http://es-search-7.fiverrdev.com:9200"
	INDEX_NAME  = "bar_books_store_2"
	QUERY_SIZE = 100
)

type ElasticDb struct {
	client *elastic.Client
	isInit bool
}

func (e *ElasticDb) Init() error{
	if !e.isInit {
		client, err := elastic.NewClient(elastic.SetURL(ELASTIC_URL))
		if err != nil {
			return err
		}
		e.client = client
		ctx := context.Background()
		_, _, err = e.client.Ping(ELASTIC_URL).Do(ctx)
		if err != nil {
			return err
		}
		e.isInit = true
	}
	return nil
}

func (e *ElasticDb) Get(id string) (*models.Book, error){
	initErr := e.Init()
	if initErr != nil {
		return nil, initErr
	}
	ctx := context.Background()
	get, getErr := e.client.Get().Index(INDEX_NAME).Id(id).Do(ctx)
	if getErr != nil {
		return nil, getErr
	}
	var book models.Book
	unmarshalErr := json.Unmarshal(get.Source, &book)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &book, nil
}

func (e *ElasticDb) Put(book models.Book) (string, error) {
	initErr := e.Init()
	if initErr != nil {
		return "", initErr
	}
	ctx := context.Background()
	bookJson, marshalErr := json.Marshal(book)
	if marshalErr != nil {
		return "", marshalErr
	}
	bj := string(bookJson)
	put, putErr := e.client.Index().
		Index(INDEX_NAME).
		BodyJson(bj).
		Do(ctx)
	if putErr != nil {
		return "", putErr
	}
	return put.Id, nil

}

func (e *ElasticDb) Post(shortBook models.ShortBook) (*models.Book, error){
	initErr := e.Init()
	if initErr != nil {
		return nil, initErr
	}
	ctx := context.Background()
	_, updateErr := e.client.Update().Index(INDEX_NAME).Id(shortBook.Id).Doc(map[string]interface{}{"title": shortBook.Title}).Do(ctx)
	if updateErr != nil {
		return nil, updateErr
	}
	return e.Get(shortBook.Id)
}

func (e *ElasticDb) Delete(id string) error {
	initErr := e.Init()
	if initErr != nil {
		return initErr
	}
	ctx := context.Background()
	_, deleteErr := e.client.Delete().Index(INDEX_NAME).Id(id).Do(ctx)
	return deleteErr
}

func (e *ElasticDb) queryBuilder(bookQuery models.BookQuery) (*elastic.BoolQuery, error){
	query := elastic.NewBoolQuery()
	if bookQuery.Title != "" {
		titleQuery := elastic.NewMatchQuery("title", bookQuery.Title)
		query.Must(titleQuery)
	}
	if bookQuery.AuthorName != "" {
		authorQuery := elastic.NewMatchQuery("author_name", bookQuery.AuthorName)
		query.Must(authorQuery)
	}
	if bookQuery.PriceRange != "" {
		rangeSlice := strings.Split(bookQuery.PriceRange, ",")
		low, lowErr := strconv.ParseFloat(rangeSlice[0], 64)
		high, highErr := strconv.ParseFloat(rangeSlice[1], 64)
		if lowErr != nil || highErr != nil {
			return nil, lowErr
		}
		rangeQuery := elastic.NewRangeQuery("price").Gte(low).Lte(high)
		query.Must(rangeQuery)
	}
	return query, nil
}

func (e *ElasticDb) Search(bookQuery models.BookQuery) ([]models.Book, error) {
	initErr := e.Init()
	if initErr != nil {
		return nil, initErr
	}
	ctx := context.Background()
	var allBooks []models.Book
	query, builderErr := e.queryBuilder(bookQuery)
	if builderErr != nil {
		return nil, builderErr
	}
	searchResult, searchErr := e.client.Search().Index(INDEX_NAME).Query(query).Size(QUERY_SIZE).Do(ctx)
	if searchErr != nil {
		return nil, searchErr
	}
	if searchResult.Hits.TotalHits.Value > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var book models.Book
			unmarshalErr := json.Unmarshal(hit.Source, &book)
			if unmarshalErr != nil {
				return nil, unmarshalErr
			}
			allBooks = append(allBooks, book)
		}
	}
	return allBooks, nil

}


func (e *ElasticDb) GetStore() (models.Store, error) {
	initErr := e.Init()
	if initErr != nil {
		return models.Store{}, initErr
	}
	ctx := context.Background()
	distinctAuthorsAgg := elastic.NewCardinalityAggregation().Field("author_name.keyword")
	allBooksSearchResults, searchErr := e.client.Search().Index(INDEX_NAME).Aggregation("bar", distinctAuthorsAgg).Size(0).Do(ctx)
	if searchErr != nil {
		return models.Store{}, searchErr
	}
	if allBooksSearchResults.Aggregations["bar"] == nil {
		// Should handle this - use found
	}

	var numOfDistinctAuthors struct{
		Count int64 `json:"value"`
	}
	unmarshalErr := json.Unmarshal(allBooksSearchResults.Aggregations["bar"], &numOfDistinctAuthors)
	if unmarshalErr != nil {
		return models.Store{}, unmarshalErr
	}
	return models.Store{
		Books: allBooksSearchResults.Hits.TotalHits.Value,
		Authors: numOfDistinctAuthors.Count,
	}, nil
}
