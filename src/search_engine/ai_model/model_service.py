import os
import logging
from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer
from rank_bm25 import BM25Okapi
import numpy as np
import pickle
from dotenv import load_dotenv
import pathlib
from fastembed import LateInteractionTextEmbedding

env_path = pathlib.Path(__file__).parents[1] / 'config' / 'config.env'
load_dotenv(env_path)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Конфигурация batch-обработки
BATCH_SIZE = int(os.getenv('BATCH_SIZE', 32))  # Размер батча по умолчанию
MAX_BATCH_SIZE = int(os.getenv('MAX_BATCH_SIZE', 256))  # Максимальный размер батча

# Загружаем dense модель
logger.info("Loading dense model...")
dense_model = SentenceTransformer(
    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2",
    device="cpu"
)
logger.info("Dense model loaded")

# Для BM25 нам нужен индекс. Создаем его из корпуса
logger.info("Loading BM25 index...")
bm25_index_path = os.getenv("bm25_index_path", "./bm25_index.pkl")

if os.path.exists(bm25_index_path):
    with open(bm25_index_path, 'rb') as f:
        bm25_index = pickle.load(f)
    logger.info("BM25 index loaded from cache")
else:
    # Если индекса нет, создаем заглушку
    logger.warning("BM25 index not found, creating dummy index")
    dummy_corpus = ["манометр избыточного давления", "манометр показывающий", "термометр", "Волтьометр","амперметр", "температура", "темпе"]
    tokenized_corpus = [doc.split() for doc in dummy_corpus]
    bm25_index = BM25Okapi(tokenized_corpus)
    logger.info("Dummy BM25 index created")

logger.info("Loading late interaction model...")
late_model = SentenceTransformer("answerdotai/answerai-colbert-small-v1")

logger.info("Late interaction model loaded")

def batch_generator(texts, batch_size):
    """Генератор для разбиения текстов на батчи"""
    for i in range(0, len(texts), batch_size):
        yield texts[i:i + batch_size]

@app.route('/embed/dense', methods=['POST'])
def embed_dense():
    try:
        data = request.json
        
        # Поддержка как одиночного текста, так и batch
        if 'text' in data:
            texts = [data['text']]
        elif 'texts' in data:
            texts = data['texts']
        else:
            return jsonify({'error': 'No text or texts provided'}), 400
        
        # Валидация размера batch
        if len(texts) > MAX_BATCH_SIZE:
            return jsonify({
                'error': f'Batch size exceeds maximum allowed ({MAX_BATCH_SIZE})'
            }), 400
        
        # Batch encoding с прогрессивной обработкой
        embeddings = dense_model.encode(
            texts,
            batch_size=BATCH_SIZE,
            normalize_embeddings=True,
            show_progress_bar=False
        )
        
        # Возвращаем одиночный embedding если был один текст, иначе массив
        if 'text' in data:
            return jsonify({'embedding': embeddings[0].tolist()})
        else:
            return jsonify({'embeddings': embeddings.tolist()})
            
    except Exception as e:
        logger.error(f"Error in dense embedding: {e}")
        return jsonify({'error': str(e)}), 500

@app.route('/embed/sparse', methods=['POST'])
def embed_sparse():
    try:
        data = request.json
        
        # Поддержка как одиночного текста, так и batch
        if 'text' in data:
            texts = [data['text']]
        elif 'texts' in data:
            texts = data['texts']
        else:
            return jsonify({'error': 'No text or texts provided'}), 400
        
        # Валидация размера batch
        if len(texts) > MAX_BATCH_SIZE:
            return jsonify({
                'error': f'Batch size exceeds maximum allowed ({MAX_BATCH_SIZE})'
            }), 400
        
        # Обработка batch для sparse embeddings
        sparse_embeddings = []
        for text in texts:
            tokenized_query = text.lower().split()
            
            # Создаем sparse представление на основе частоты токенов
            word_counts = {}
            for token in tokenized_query:
                word_counts[token] = word_counts.get(token, 0) + 1
            
            # Создаем индексы (хэши слов) и значения (TF-IDF weights)
            indices = [hash(word) % 10000 for word in word_counts.keys()]
            values = [count / len(tokenized_query) for count in word_counts.values()]
            
            sparse_embeddings.append({
                'indices': indices,
                'values': values
            })
        
        if 'text' in data:
            return jsonify(sparse_embeddings[0])
        else:
            return jsonify({'sparse_embeddings': sparse_embeddings})
            
    except Exception as e:
        logger.error(f"Error in sparse embedding: {e}")
        return jsonify({'error': str(e)}), 500

@app.route('/embed/late', methods=['POST'])
def embed_late():
    try:
        data = request.json
        
        # Поддержка как одиночного текста, так и batch
        if 'text' in data:
            texts = [data['text']]
        elif 'texts' in data:
            texts = data['texts']
        else:
            return jsonify({'error': 'No text or texts provided'}), 400
        
        # Валидация размера batch
        if len(texts) > MAX_BATCH_SIZE:
            return jsonify({
                'error': f'Batch size exceeds maximum allowed ({MAX_BATCH_SIZE})'
            }), 400
        
        multi_vectors = []
        for text in texts:
            token_embeddings = late_model.encode(
                text, 
                output_value='token_embeddings',
                normalize_embeddings=True
            )
            multi_vectors.append(token_embeddings.tolist())
        
        if 'text' in data:
            return jsonify({'multi_vector': multi_vectors[0]})
        else:
            return jsonify({'multi_vectors': multi_vectors})
            
    except Exception as e:
        logger.error(f"Error in late embedding: {e}")
        return jsonify({'error': str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        'status': 'ok',
        'models': 'loaded',
        'batch_size': BATCH_SIZE,
        'max_batch_size': MAX_BATCH_SIZE
    })

# создаем индекс для BM25, елсли его нет, иначе обновляем
@app.route('/update_index', methods=["POST"])
def update_index():
    global bm25_index_path, bm25_index
    data = request.json
    new_corpus = data['miinstance_name']
    
    # Проверка, что новые данные не пустые
    if not new_corpus:
        return {"error": "Нет данных для добавления"}, 400
    
    # Если это строка - превращаем в список
    if isinstance(new_corpus, str):
        new_corpus = [new_corpus]
    
    # Если индекс не существует - создаём с нуля
    if not os.path.exists(bm25_index_path):
        tokenized = [doc.lower().split() for doc in new_corpus]
        bm25 = BM25Okapi(tokenized)
        
        with open(bm25_index_path, 'wb') as f:
            pickle.dump({'index': bm25, 'corpus': new_corpus}, f)
        bm25_index = bm25
        
        return {
            "message": "Индекс создан", 
            "documents_count": len(new_corpus)
        }, 200
    
    # Индекс существует - обновляем
    else:
        bm25_index_path = 'bm25_index.pkl'
        with open(bm25_index_path, 'rb') as f:
            saved = pickle.load(f)
            existing_corpus = saved.get('corpus', [])
        
        # Объединяем старые и новые документы
        full_corpus = existing_corpus + new_corpus
        
        # Пересоздаём индекс с нуля
        tokenized_corpus = [doc.lower().split() for doc in full_corpus]
        bm25 = BM25Okapi(tokenized_corpus)
        
        # Сохраняем
        with open(bm25_index_path, 'wb') as f:
            pickle.dump({
                'index': bm25, 
                'corpus': full_corpus
            }, f)
        
        bm25_index = bm25
        
        return {
            "message": "Индекс обновлён",
            "previous_count": len(existing_corpus),
            "added_count": len(new_corpus),
            "total_count": len(full_corpus)
        }, 200

if __name__ == '__main__':
    port = int(os.getenv('PORT', 8000))
    host = os.getenv('HOST', '0.0.0.0')
    logger.info(f"Starting embedding service on {host}:{port}")
    logger.info(f"Batch size: {BATCH_SIZE}, Max batch size: {MAX_BATCH_SIZE}")
    app.run(host=host, port=port, debug=True, threaded=True)