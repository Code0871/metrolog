import os
import pathlib
from dotenv import load_dotenv

# загружаем .env один раз при импорте
env_path = pathlib.Path(__file__).parents[2] / 'config' / 'config.env'
load_dotenv(env_path)

class Config:
    # model settings
    models_cache = os.getenv('models_cache', '/app/models_cache')

    dense_model_name = os.getenv('dense_model_name', 'sentence-transformers/paraphrase-multilingual-mini_lm-l12-v2')
    late_model_name = os.getenv('late_model_name', 'answerdotai/answerai-colbert-small-v1')
    device = os.getenv('device', 'cpu')
    
    # batch settings
    batch_size = int(os.getenv('batch_size', 32))
    max_batch_size = int(os.getenv('max_batch_size', 256))
    
    # bm25 settings
    bm25_index_path = os.getenv('bm25_index_path', './bm25_index.pkl')
    
    # server settings
    port = int(os.getenv('port', 8000))
    host = os.getenv('host', '0.0.0.0')
    debug = os.getenv('debug', 'true').lower() == 'true'