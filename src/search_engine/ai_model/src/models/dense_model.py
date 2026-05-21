from sentence_transformers import SentenceTransformer
from pathlib import Path
from typing import List
import os

class DenseEmbeddingModel:
    def __init__(self, model_name: str, cache_dir: str = "./models_cache"):
        self.model_name = model_name
        self.cache_dir = Path(cache_dir)
        
        model_path = self._get_model_path()
        
        # Загружаем модель из локального пути
        if model_path.exists():
            print(f"Loading model from local path: {model_path}")
            self.model = SentenceTransformer(str(model_path))
        else:
            print(f"Local model not found, downloading: {model_name}")
            self.model = SentenceTransformer(
                model_name, 
                cache_folder=str(self.cache_dir)
            )
    
    def _get_model_path(self) -> Path:
        """Получить путь к скачанной модели"""
        # snapshot_download сохраняет в cache_dir/models--org--model_name
        model_dir_name = "models--" + self.model_name.replace("/", "--")
        return self.cache_dir / model_dir_name / "snapshots" / "*"
    
    def encode(self, texts: List[str]) -> List[List[float]]:
        """Получить эмбеддинги"""
        embeddings = self.model.encode(texts, convert_to_numpy=True)
        return embeddings.tolist()