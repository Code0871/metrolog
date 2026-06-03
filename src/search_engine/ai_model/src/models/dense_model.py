from sentence_transformers import SentenceTransformer
from pathlib import Path
from typing import List
import os

class DenseEmbeddingModel:
    def __init__(self, model_name: str, models_dir: str = "./models_dir/dense"):
        self.model_name = model_name

        self.models_dir = Path(__file__).parents[2] / "models_dir" / "dense"
        
        if self.models_dir.exists() and any(self.models_dir.iterdir()):
            print(f"Loading model from local path: {self.models_dir}")
            self.model = SentenceTransformer(str(self.models_dir))
        else:
            print(f"Local model not found at {self.models_dir}, downloading: {model_name}")
            self.model = SentenceTransformer(
                model_name, 
                cache_folder=str(self.models_dir)
            )
    
    def encode(self, texts: List[str]) -> List[List[float]]:
        """Получить эмбеддинги"""
        embeddings = self.model.encode(texts, convert_to_numpy=True)
        return embeddings.tolist()