
import logging
from qdrant_client import QdrantClient
from qdrant_client.http import models

logger = logging.getLogger(__name__)

class qdrant_rep:
    def __init__(self, client: QdrantClient, collection_name: str):
        self.client = client
        self.collection_name = collection_name

    def create_collection(self, vector_size: int = 384, distance: str = "COSINE") -> bool:
        
        if self.client.collection_exists(self.collection_name):
            logger.info(f"collection '{self.collection_name}' is already exists")
            return False
        
        logger.info(f"creating collection '{self.collection_name}'")
        self.client.create_collection(
            collection_name = self.collection_name,
            vectors_config = models.VectorParams(
                size = vector_size,
                distance = models.Distance[distance.upper()],
            )
        )

        logger.info(f"collection '{self.collection_name}' created")
        return True