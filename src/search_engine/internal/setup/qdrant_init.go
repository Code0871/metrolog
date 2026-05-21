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

// инициализация коллекции для гибридного поиска
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
		//client.DeleteCollection(context.Background(), collection_name)
		return
	}

	m := uint64(16)
	ef_construct := uint64(100)
	full_scan_threshold := uint64(10000)
	on_disk := false
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
		VectorsConfig: qdrant.NewVectorsConfigMap(
			map[string]*qdrant.VectorParams{
				"dense": {
					Size:     vec_size,
					Distance: distance_type,
					HnswConfig: &qdrant.HnswConfigDiff{
						M:                 &m,
						EfConstruct:       &ef_construct,
						FullScanThreshold: &full_scan_threshold,
						OnDisk:            &on_disk,
						PayloadM:          &payload_m,
					},
				},
				"multi": {
					Size:     384,
					Distance: qdrant.Distance_Cosine,
					MultivectorConfig: &qdrant.MultiVectorConfig{
						Comparator: qdrant.MultiVectorComparator_MaxSim,
					},
					HnswConfig: &qdrant.HnswConfigDiff{M: qdrant.PtrOf(uint64(0))},
				},
			}),
		SparseVectorsConfig: qdrant.NewSparseVectorsConfig(
			map[string]*qdrant.SparseVectorParams{
				"sparse": {Modifier: qdrant.Modifier_Idf.Enum()},
			},
		),
	})
	if err != nil {
		panic(err)
	}
}
