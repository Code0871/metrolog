package setup

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

// инициализация клиента Qdrant
func MustInitQdrantСlient(host string, port int) *qdrant.Client {
	cli, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		panic(fmt.Sprintf("qdrant init failed: %v", err))
	}

	_, err = cli.HealthCheck(context.Background())
	if err != nil {
		panic(fmt.Sprintf("qdrant check health failed: %v", err))
	}

	return cli
}

// TODO: можно сделать передачу параметров для построения индексов читая их настройки из конфига
// инициализация коллекции
func MustInitQdrantCollection(client *qdrant.Client, collection_name string, vec_size uint64, distance_type qdrant.Distance) {

	if collections, err := client.ListCollections(context.Background()); err == nil {
		for _, val := range collections {
			if val == collection_name {
				fmt.Println("Collection already exists")
				return
			}
		}

		fmt.Println("Collection not found")
	}

	m := uint64(16)
	efconstruction := uint64(100)
	fullscanthreshhold := uint64(10000)
	ondisk := true
	payloadm := uint64(100)

	err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collection_name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vec_size,
			Distance: distance_type,
			HnswConfig: &qdrant.HnswConfigDiff{
				M:                 &m,
				EfConstruct:       &efconstruction,
				FullScanThreshold: &fullscanthreshhold,
				OnDisk:            &ondisk,
				PayloadM:          &payloadm,
			},
		}),
	})
	if err != nil {
		panic(err)
	}
}
