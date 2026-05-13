package setup

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func MustInitQdrantСlient(ctx context.Context, host string, port int) *qdrant.Client {
	cli, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		panic(fmt.Sprintf("qdrant init failed: %v", err))
	}

	_, err = cli.HealthCheck(ctx)
	if err != nil {
		panic(fmt.Sprintf("qdrant check health failed: %v", err))
	}

	return cli
}

func MustInitQdrantCollection(client *qdrant.Client, collection_name string, vec_size uint64, distance_type qdrant.Distance) {
	err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collection_name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vec_size,
			Distance: distance_type,
		}),
	})
	if err != nil {
		panic(err)
	}
}
