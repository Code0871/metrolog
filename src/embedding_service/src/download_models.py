import os
import logging
from pathlib import Path
from huggingface_hub import snapshot_download
from concurrent.futures import ThreadPoolExecutor, as_completed

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def download_model(repo_id: str, model_dir: str = "./models_dir") -> str:
    """Скачать любую модель с HuggingFace одним вызовом"""
    logger.info(f"Downloading {repo_id}...")
    
    return snapshot_download(
        repo_id=repo_id,
        local_dir=model_dir,
        local_dir_use_symlinks=False,
        resume_download=True
    )


def main():
    """Скачать все модели из конфига"""
    from dotenv import load_dotenv
    
    if not os.path.exists('models_dir'):
        os.mkdir('models_dir')
        os.mkdir('models_dir/dense')
        os.mkdir('models_dir/sparse')
        os.mkdir('models_dir/late')
    
    print(f"Проверка директории с моделями: {os.listdir('models_dir')}")
    if os.listdir('models_dir') == ['sparse', 'dense', 'late']:
        logger.info("Проверяем наличие моделей")
        if len(os.listdir('models_dir/sparse')) != 0 and len(os.listdir('models_dir/dense')) != 0  and len(os.listdir('models_dir/late')) != 0 :
            logger.info("Модели уже скачаны")
            return

    # Грузим конфиг
    env_path = os.getenv("CONFIG_PATH", "../config/config.env")
    load_dotenv(env_path)
    
    # Список моделей для скачивания (можно дополнить)
    models = {
        os.getenv("dense_embedding_model", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"): "dense",
        os.getenv("sparse_model", "Qdrant/bm25"): "sparse",
        os.getenv("late_interacction_embeding_model", "answerdotai/answerai-colbert-small-v1"): "late",
    }

    # Фильтруем пустые значения и скачиваем
    models_dir = os.getenv("models_dir", "./models_dir")
    
    with ThreadPoolExecutor() as executor:
        futures = {}
        for model, subfolder in models.items():
            if not model:
                continue
                
            target_dir = os.path.join(models_dir, subfolder)
            os.makedirs(target_dir, exist_ok=True)
            
            # Асинхронный запуск (не блокирует цикл)
            future = executor.submit(download_model, model, target_dir)
            futures[future] = model
        
        # Собираем результаты по мере завершения
        for future in as_completed(futures):
            model = futures[future]
            try:
                path = future.result()
                logger.info(f"Модель {model} скачана в папку: {path}")
            except Exception as e:
                logger.error(f"Ошибка скачивания модели {model}: {e}")

    if len(os.listdir('models_dir/dense')) == 0:
        logger.info(f"Пробуем уставновить модель {models['dense']} повторно")
        download_model(models['dense'], 'models_dir/dense')
    if len(os.listdir('models_dir/sparse')) == 0:
        logger.info(f"Пробюуем уставновить модель {models['sparse']} повторно")
        download_model(models['sparse'], 'models_dir/sparse')
    if len(os.listdir('models_dir/late')) == 0:
        download_model(models['late'], 'models_dir/late')

if __name__ == "__main__":
    main()