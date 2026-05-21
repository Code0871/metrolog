# model_service.py
import os
from pathlib import Path
from src.services.embedding_service import EmbeddingService
from typing import List


class ModelService:
    def __init__(self):
        self.embedding_service = EmbeddingService()
    
    def encode(self, texts: List[str], method: str = "dense"):
        if method == "dense":
            return self.embedding_service.encode_dense(texts)
        elif method == "sparse":
            return self.embedding_service.encode_sparse(texts)
        elif method == "late":
            return self.embedding_service.encode_late(texts)
        elif method == "all":
            return self.embedding_service.encode_all(texts)
        else:
            raise ValueError(f"Unknown method: {method}")
    
    def search(self, query: str, documents: List[str], method: str = "dense", top_k: int = 5):
        if method == "late":
            return self.embedding_service.search_late(query, documents, top_k)
        else:
            query_emb = self.encode([query], method)[0]
            doc_embs = self.encode(documents, method)
            
            import numpy as np
            query_arr = np.array(query_emb)
            scores = []
            for i, doc_emb in enumerate(doc_embs):
                doc_arr = np.array(doc_emb)
                similarity = np.dot(query_arr, doc_arr) / (
                    np.linalg.norm(query_arr) * np.linalg.norm(doc_arr)
                )
                scores.append({"index": i, "score": float(similarity)})
            
            scores.sort(key=lambda x: x["score"], reverse=True)
            return scores[:top_k]