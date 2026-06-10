from pathlib import Path
from typing import List, Optional
import pickle
import os

from src.services.bm25_indexer import get_bm25_indexer, bm25_indexer

class SparseModel:
    def __init__(self, model_name: str = "Qdrant/bm25", models_dir: str = "./models_dir/sparse"):
        self.model_name = model_name
        self.cache_dir = Path(models_dir)
        
        # Используем уже готовый индекс из сервисного слоя
        self.bm25 = bm25_indexer
        self.is_ready = self.bm25 is not None
        
        # Загружаем корпус (список документов)
        self.corpus = self.load_corpus()
        
        if self.is_ready:
            print(f"SparseModel готов, документов: {len(self.corpus) if self.corpus else 'неизвестно'}")
        else:
            print("SparseModel: индекс не загружен")
    
    def load_corpus(self) -> List[str]:
        index_path = 'bm25_index.pkl'
        
        if os.path.exists(index_path):
            with open(index_path, 'rb') as f:
                data = pickle.load(f)
                return data.get('corpus', [])
        return []
    
    def search(self, query: str, top_k: int = 10) -> List[str]:
        if not self.is_ready:
            raise ValueError("BM25 индекс не загружен. Запустите create_bm25_simple() сначала")
        
        # Токенизация запроса (как при создании индекса)
        tokenized_query = query.lower().split()
        
        # Получение топ-N документов
        top_docs = self.bm25.get_top_n(tokenized_query, self.corpus, n=top_k)
        
        return top_docs
    
    def get_scores(self, query: str) -> List[float]:
        if not self.is_ready:
            raise ValueError("BM25 индекс не загружен")
        
        tokenized_query = query.lower().split()
        scores = self.bm25.get_scores(tokenized_query)
        return scores.tolist()
    
    def search_with_scores(self, query: str, top_k: int = 10) -> List[tuple]:
        if not self.is_ready:
            raise ValueError("BM25 индекс не загружен")
        
        tokenized_query = query.lower().split()
        scores = self.bm25.get_scores(tokenized_query)
        
        # Сортируем документы по убыванию оценки
        scored_docs = sorted(zip(self.corpus, scores), key=lambda x: x[1], reverse=True)
        
        return scored_docs[:top_k]