from flask import request, jsonify
from src.services.embedding_service import EmbeddingService
import numpy as np

# Используем синглтон
embedding_service = EmbeddingService()

def convert_to_serializable(obj):
    """Рекурсивно преобразует numpy типы в Python типы"""
    if isinstance(obj, np.ndarray):
        return obj.tolist()
    elif isinstance(obj, np.integer):
        return int(obj)
    elif isinstance(obj, np.floating):
        return float(obj)
    elif isinstance(obj, dict):
        return {k: convert_to_serializable(v) for k, v in obj.items()}
    elif isinstance(obj, (list, tuple)):
        return [convert_to_serializable(item) for item in obj]
    return obj

def register_late_routes(app):
    @app.route('/embed/late', methods=['GET'])
    def embed_late():
        """Получить мультивекторы (token-level embeddings)"""
        try:
            data = request.json
            
            if 'text' in data:
                texts = [data['text']]
                is_single = True
            elif 'texts' in data:
                texts = data['texts']
                is_single = False
            else:
                return jsonify({'error': 'No text or texts provided'}), 400
            
            # Получаем мультивекторы
            multi_vectors = convert_to_serializable(embedding_service.encode_late(texts))
            print(multi_vectors)
            if is_single:
                return jsonify({'multi_vector': multi_vectors[0]})
            else:
                return jsonify({'multi_vectors': multi_vectors})
                
        except Exception as e:
            return jsonify({'error': str(e)}), 500