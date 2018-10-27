package search

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gonzalesraul/meow/schema"
	"github.com/olivere/elastic"
)

//ElasticRepository keeps the elasticsearch client reference
type ElasticRepository struct {
	client *elastic.Client
}

//NewElastic created a new connection with ElasticSearch and return it
func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

//Close the connection with ElasticSearch
func (r *ElasticRepository) Close() {
	//TODO
}

//InsertMeow performs an insert on ElasticSearch
func (r *ElasticRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
	_, err := r.client.Index().
		Index("meow").
		Type("meow").
		Id(meow.ID).
		BodyJson(meow).
		Refresh("wait_for").
		Do(ctx)
	return err
}

//InquiryMeow return the meows finds eventually for the query provided
func (r *ElasticRepository) InquiryMeow(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error) {
	result, err := r.client.Search().
		Index("meow").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		return nil, err
	}
	meows := []schema.Meow{}
	for _, hit := range result.Hits.Hits {
		var meow schema.Meow
		if err = json.Unmarshal(*hit.Source, &meow); err != nil {
			log.Println(err)
		}
		meows = append(meows, meow)
	}
	return meows, nil
}
