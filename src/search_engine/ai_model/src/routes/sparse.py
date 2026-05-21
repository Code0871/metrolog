from flask import request, jsonify
from src.services.embedding_service import EmbeddingService

# Создаем один экземпляр сервиса
embedding_service = EmbeddingService()

def register_sparse_routes(app):
    @app.route('/embed/sparse', methods=['POST'])
    def embed_sparse():
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
            
            sparse_embeddings = embedding_service.process_sparse(texts)
            
            if is_single:
                return jsonify(sparse_embeddings[0])
            else:
                return jsonify({'sparse_embeddings': sparse_embeddings})
                
        except ValueError as e:
            return jsonify({'error': str(e)}), 400
        except Exception as e:
            return jsonify({'error': str(e)}), 500