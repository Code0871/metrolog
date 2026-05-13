package setup

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

func InitQdrantСlient(ctx context.Context, host string, port int) (*qdrant.Client, error) {
	cli, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	_, err = cli.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func InitQdrantCollections() {

}
