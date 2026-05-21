from flask import request, jsonify
from src.services.embedding_service import EmbeddingService

embedding_service = EmbeddingService()

def register_dense_routes(app):
    @app.route('/embed/dense', methods=['POST'])
    def embed_dense():
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
            
            embeddings = embedding_service.process_dense(texts)
            
            if is_single:
                return jsonify({'embedding': embeddings[0].tolist()})
            else:
                return jsonify({'embeddings': embeddings.tolist()})
                
        except ValueError as e:
            return jsonify({'error': str(e)}), 400
        except Exception as e:
            return jsonify({'error': str(e)}), 500