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

// инициализация коллекции
func MustInitQdrantCollection(client *qdrant.Client, collection_name string, vec_size uint64, distance_type qdrant.Distance) {

	// Проверяем, что клиент существует
	if client == nil {
		panic("qdrant client is not exists")
	}

	// Проверяем, что коллецкция с заданным именем не существует
	collection, err := client.CollectionExists(context.Background(), collection_name)
	if err != nil {
		panic(err)
	}
	if collection {
		fmt.Printf("Collection '%s' already exists\n", collection_name)
		return
	}

	// параметры Hnsw
	m := uint64(16)
	ef_construction := uint64(100)
	full_scan_threshold := uint64(10000)
	on_disk := true
	payload_m := uint64(100)

	// параметры вакуумного оптимизатора
	delete_threshold := float64(0.2)
	vacuum_min_vector_number := uint64(500)

	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collection_name,
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			DeletedThreshold:      &delete_threshold,
			VacuumMinVectorNumber: &vacuum_min_vector_number,
		},
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vec_size,
			Distance: distance_type,
			HnswConfig: &qdrant.HnswConfigDiff{
				M:                 &m,
				EfConstruct:       &ef_construction,
				FullScanThreshold: &full_scan_threshold,
				OnDisk:            &on_disk,
				PayloadM:          &payload_m,
			},
		}),
	})
	if err != nil {
		panic(err)
	}
}
