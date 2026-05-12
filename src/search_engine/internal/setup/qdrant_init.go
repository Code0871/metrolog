package setup

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

func InitQdrantСlient(ctx context.Context, host string, port int) (*qdrant.Client, error) {
	cli, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	return cli, nil
}

func InitQdrantCollections() {
	
}