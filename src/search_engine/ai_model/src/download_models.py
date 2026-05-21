# src/models/download_models.py
import os
import logging
from pathlib import Path
from huggingface_hub import snapshot_download

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def download_model(repo_id: str, cache_dir: str = "./models_cache") -> str:
    """Скачать любую модель с HuggingFace одним вызовом"""
    logger.info(f"Downloading {repo_id}...")
    
    return snapshot_download(
        repo_id=repo_id,
        cache_dir=cache_dir,
        local_dir_use_symlinks=False,
        resume_download=True
    )


def main():
    """Скачать все модели из конфига"""
    from dotenv import load_dotenv
    
    # Грузим конфиг
    env_path = os.getenv("CONFIG_PATH", "../config/config.env")
    load_dotenv(env_path)
    
    # Список моделей для скачивания (можно дополнить)
    models = [
        os.getenv("dense_embedding_model", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"),
        os.getenv("sparse_model", "Qdrant/bm25"),
        os.getenv("late_interacction_embeding_model", "answerdotai/answerai-colbert-small-v1"),
    ]
    
    # Фильтруем пустые значения и скачиваем
    cache_dir = os.getenv("model_cache", "./models_cache")
    
    for model in filter(None, models):
        try:
            path = download_model(model, cache_dir)
            logger.info(f"Downloaded to: {path}")
        except Exception as e:
            logger.error(f"Failed to download {model}: {e}")


if __name__ == "__main__":
    main()