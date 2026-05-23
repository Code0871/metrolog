import os
import numpy as np
from typing import List, Dict, Optional
from pathlib import Path
from src.models.sparse_model import SparseModel
from src.services.bm25_indexer import create_bm25_simple, get_bm25_indexer
from dotenv import load_dotenv

# Правильный путь к конфигу
env_path = Path(__file__).parents[3] / 'config' / 'config.env'
print(f"Config path: {env_path}")
load_dotenv(env_path)

class EmbeddingService:
    _instance = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialized = False
        return cls._instance
    
    def __init__(self):
        if self._initialized:
            return
        
        self.cache_dir = os.getenv("model_cache", "./models_cache")

        # Sparse модель
        create_bm25_simple()
        bm25 = get_bm25_indexer()
        self.sparse_model = SparseModel(
            model_name=os.getenv("sparse_model", "Qdrant/bm25"),
            cache_dir=self.cache_dir
        )
        
        # Dense модель
        from sentence_transformers import SentenceTransformer
        dense_name = os.getenv("dense_embedding_model", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")
        print(f"Loading dense model: {dense_name}")
        self.dense_model = SentenceTransformer(dense_name, cache_folder=self.cache_dir)
        print("Dense model loaded!")
        
        # Late модель
        late_name = (os.getenv("late_interaction_embedding_model") or 
                    os.getenv("late_interacction_embeding_model"))
        
        print(f"Late model name from env: {late_name}")
        
        if late_name:
            from sentence_transformers import SentenceTransformer
            print(f"Loading late model: {late_name}")
            self.late_model = SentenceTransformer(late_name, cache_folder=self.cache_dir)
            print("Late model loaded!")
        else:
            print("Late model not configured in .env!")
            self.late_model = None
        
        self._initialized = True
    
    def process_dense(self, texts: List[str]) -> np.ndarray:
        if self.dense_model is None:
            raise ValueError("Dense model not initialized")
        return self.dense_model.encode(texts, normalize_embeddings=True)
    
    def process_sparse(self, texts: List[str]) -> List[Dict]:
        return self.sparse_model.encode(texts)
    
    def encode_late(self, texts: List[str]) -> List[List[List[float]]]:
        if self.late_model is None:
            raise ValueError("Late model not initialized")
        
        multi_vectors = []
        for text in texts:
            token_embeddings = self.late_model.encode(
                text,
                output_value='token_embeddings',
                normalize_embeddings=True
            )
            multi_vectors.append(token_embeddings.tolist())
        
        return multi_vectors
    
    def update_bm25_index(self, corpus):
        return self.sparse_model.update_index(corpus)