# src/routes/late.py
from flask import request, jsonify
from src.services.embedding_service import EmbeddingService

# Используем синглтон
embedding_service = EmbeddingService()

def register_late_routes(app):
    @app.route('/embed/late', methods=['POST'])
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
            multi_vectors = embedding_service.encode_late(texts)
            
            if is_single:
                return jsonify({'multi_vector': multi_vectors[0]})
            else:
                return jsonify({'multi_vectors': multi_vectors})
                
        except Exception as e:
            return jsonify({'error': str(e)}), 500