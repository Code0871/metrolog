import os
import logging
from qdrant_client import QdrantClient
from internal.repository.qdrant_rep import QdrantRepository

logger = logging.getLogger(__name__)

def main():
    logger.info("initial qdrant")
    
    qdrant_url = os.getenv("QDRANT_URL", "http://localhost:6333")
    collection_name = os.getenv("QDRANT_COLLECTION", "miinstance")
    vector_size = int(os.getenv("VECTOR_SIZE", "384"))
    distance = os.getenv("DISTANCE", "COSINE")

    client = QdrantClient(url=qdrant_url, timeout = 10.0)

    repo = QdrantRepository(client=client, collection_name=collection_name)

    repo.create_collection(
        vector_size = vector_size,
        distance = distance
    )


if __name__ == "__main__":
    main()