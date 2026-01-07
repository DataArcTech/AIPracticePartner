package Common

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func EsClient(ctx context.Context) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		return nil, fmt.Errorf("create es client failed: %w", err)
	}
	return client, nil
}
