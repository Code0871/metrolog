# model_service.py
import os
import numpy as np
from pathlib import Path
from src.services.embedding_service import EmbeddingService
from typing import List, Dict, Any


class ModelService:
    def __init__(self):
        self.embedding_service = EmbeddingService()
    
    def encode(self, texts: List[str], method: str = "dense") -> List:
        if method == "dense":
            return self.embedding_service.process_dense(texts)
        elif method == "sparse":
            # Для sparse (BM25) возвращаем результаты поиска, не эмбеддинги
            # BM25 не имеет отдельных эмбеддингов, это поисковый метод
            raise ValueError("Sparse method doesn't return embeddings. Use search() instead.")
        elif method == "late":
            return self.embedding_service.encode_late(texts)
        else:
            raise ValueError(f"Unknown method: {method}. Use 'dense' or 'late'")
    
    def search(self, query: str, documents: List[str] = None, method: str = "dense", top_k: int = 5) -> List[Dict]:

        if method == "sparse":
            # Используем встроенный BM25 индекс
            results = self.embedding_service.process_sparse_with_scores(query, top_k)
            return [
                {"index": i, "document": doc, "score": score}
                for i, (doc, score) in enumerate(results)
            ]
        
        elif method == "dense":
            if documents is None:
                raise ValueError("Documents are required for dense search")
            
            # Получаем эмбеддинги
            query_emb = self.encode([query], method)[0]
            doc_embs = self.encode(documents, method)
            
            # Вычисляем косинусное сходство
            query_arr = np.array(query_emb) if not isinstance(query_emb, np.ndarray) else query_emb
            scores = []
            
            for i, doc_emb in enumerate(doc_embs):
                doc_arr = np.array(doc_emb) if not isinstance(doc_emb, np.ndarray) else doc_emb
                
                # нормализация векторов для косинусного сходства
                norm_query = np.linalg.norm(query_arr)
                norm_doc = np.linalg.norm(doc_arr)
                
                if norm_query == 0 or norm_doc == 0:
                    similarity = 0.0
                else:
                    similarity = np.dot(query_arr, doc_arr) / (norm_query * norm_doc)
                
                scores.append({
                    "index": i,
                    "document": documents[i],
                    "score": float(similarity)
                })
            
            scores.sort(key=lambda x: x["score"], reverse=True)
            return scores[:top_k]
        
        elif method == "late":
            if documents is None:
                raise ValueError("Documents are required for late search")
            
            # Для late нужно специальное сравнение (максимальное сходство по токенам)
            return self._search_late(query, documents, top_k)
        
        else:
            raise ValueError(f"Unknown method: {method}. Use 'dense', 'sparse', or 'late'")
    
    def _search_late(self, query: str, documents: List[str], top_k: int = 5) -> List[Dict]:
        # Получаем эмбеддинги для запроса (матрица токенов)
        query_emb = self.embedding_service.encode_late_query(query)
        # Получаем эмбеддинги для документов (список матриц)
        doc_embs = self.embedding_service.encode_late(documents)
        
        scores = []
        
        for i, doc_emb in enumerate(doc_embs):
            # Вычисляем матрицу сходства (num_query_tokens x num_doc_tokens)
            similarity_matrix = np.dot(query_emb, doc_emb.T)
            
            # Максимальное сходство для каждого токена запроса
            max_per_query_token = np.max(similarity_matrix, axis=1)
            
            # Итоговая оценка - сумма максимальных сходств
            score = np.sum(max_per_query_token)
            
            scores.append({
                "index": i,
                "document": documents[i],
                "score": float(score)
            })
        
        scores.sort(key=lambda x: x["score"], reverse=True)
        return scores[:top_k]
    
    def search_bm25(self, query: str, top_k: int = 10) -> List[Dict]:
        return self.search(query, method="sparse", top_k=top_k)
    
    def get_bm25_corpus(self) -> List[str]:
        return self.embedding_service.get_bm25_corpus()