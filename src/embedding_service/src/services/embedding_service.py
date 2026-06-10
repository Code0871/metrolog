import os
import numpy as np
from typing import List, Dict, Optional
from pathlib import Path
from src.models.sparse_model import SparseModel
from src.services.bm25_indexer import create_bm25_simple, get_bm25_indexer, bm25_indexer
from dotenv import load_dotenv

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

        # Sparse модель (BM25) - просто обертка над существующим индексом
        print("Initializing Sparse Model (BM25)...")
        create_bm25_simple()
        self.sparse_model = SparseModel()  # Автоматически подхватит индекс
        
        # Dense модель
        from sentence_transformers import SentenceTransformer
        dense_name = os.getenv("dense_embedding_model", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")
        print(f"Loading dense model: {dense_name}")
        self.dense_model = SentenceTransformer(dense_name, cache_folder=self.cache_dir)
        print("Dense model loaded!")
        
        # Late модель (ColBERT)
        late_name = os.getenv("late_interaction_embedding_model")
        print(f"Late model name from env: {late_name}")
        
        if late_name:
            from pylate import models as late_models
            from pathlib import Path
            
            # Путь к локальной модели ColBERT
            late_model_path = Path(__file__).parents[3] / "models_dir" / "late"
            
            if late_model_path.exists():
                print(f"Loading late model from local path: {late_model_path}")
                self.late_model = late_models.ColBERT(
                    model_name_or_path=str(late_model_path),
                    local_files_only=True
                )
            else:
                print(f"Local model not found at {late_model_path}, downloading from HF...")
                self.late_model = late_models.ColBERT(
                    model_name_or_path=late_name
                )
            print("Late model loaded!")
        else:
            print("Late model not configured in .env!")
            self.late_model = None
        
        self._initialized = True
    
    def process_dense(self, texts: List[str]) -> np.ndarray:
        if self.dense_model is None:
            raise ValueError("Dense model not initialized")
        return self.dense_model.encode(texts, normalize_embeddings=True)
    
    def process_sparse(self, texts: List[str], top_k: int = 10) -> List[dict]:
        results = []
        for text in texts:
            # Получаем BM25 оценки для этого текста как запроса
            scores = self.sparse_model.get_scores(text)
            
            # Конвертируем в numpy array, если это list
            if isinstance(scores, list):
                scores = np.array(scores)
            
            sparse_vector = self.convert_to_sparse_format(scores)
            results.append(sparse_vector)
        return results
    
    def process_sparse_with_scores(self, query: str, top_k: int = 10) -> List[tuple]:
        if self.sparse_model is None:
            raise ValueError("Sparse model not initialized")
        return self.sparse_model.search_with_scores(query, top_k=top_k)
    
    def encode_late(self, texts: List[str]) -> List:
        if self.late_model is None:
            raise ValueError("Late model not initialized")
        
        # encode возвращает список numpy массивов (num_tokens x 96)
        embeddings = self.late_model.encode(texts, is_query=False)
        return embeddings
    
    def encode_late_query(self, query: str) -> List:
        if self.late_model is None:
            raise ValueError("Late model not initialized")
        
        # Для запроса используем is_query=True
        return self.late_model.encode([query], is_query=True)[0]
    
    def get_bm25_corpus(self) -> List[str]:
        if self.sparse_model:
            return self.sparse_model.corpus
        return []
    
    def refresh_sparse_index(self):
        print("Refreshing sparse model index...")
        self.sparse_model.reload() if hasattr(self.sparse_model, 'reload') else None
        print("Sparse model refreshed!")
    
    def convert_to_sparse_format(self, scores) -> dict:

        # Убеждаемся, что scores - это numpy array
        if isinstance(scores, list):
            scores = np.array(scores)
        
        # Находим ненулевые значения
        non_zero_indices = np.where(scores > 0)[0]
        non_zero_values = scores[non_zero_indices]
        
        return {
            "indices": non_zero_indices.tolist(),
            "values": non_zero_values.tolist()
        }

# Синглтон экземпляр
embedding_service = EmbeddingService()