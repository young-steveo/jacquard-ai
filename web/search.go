package web

import (
	"context"
	"net/http"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type Searcher struct {
	engineID string
	service  *customsearch.Service
}

func NewSearcher(ctx context.Context, apiKey string, engineID string) (*Searcher, error) {
	service, err := customsearch.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Searcher{
		engineID: engineID,
		service:  service,
	}, nil
}

func (s *Searcher) Search(query string) ([]*customsearch.Result, error) {
	response, err := s.service.Cse.List().Cx(s.engineID).Q(query).Do()
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (s *Searcher) TestConnection() error {
	_, err := s.Search("is google down?")

	if err != nil {
		if gErr, ok := err.(*googleapi.Error); !ok || gErr.Code != http.StatusNotModified {
			return err
		}
	}
	return nil
}
