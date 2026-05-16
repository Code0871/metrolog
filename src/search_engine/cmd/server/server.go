package main

import (
	_ "context"
	"fmt"
	conf "search_engine/config"
	setup "search_engine/internal/setup"
)

func main() {
	fmt.Println("Load Configs")
	conf.MustLoadConfig()

	fmt.Println("Init Qdrant client")
	client := setup.MustInitQdrantСlient(conf.MustLoadConfig().QdrantConfigs.QdrantHost, conf.MustLoadConfig().QdrantConfigs.QdrantPort)

	fmt.Println("Init Qdrant collection")
	collection_name := conf.MustLoadConfig().CollectionConfig.CollectionName
	vector_size := uint64(conf.MustLoadConfig().CollectionConfig.QdrantVectorSize)
	distance_type := conf.MustLoadConfig().CollectionConfig.QdrantDistanceType

	setup.MustInitQdrantCollection(client, collection_name, vector_size, distance_type)
}
