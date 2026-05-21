from flask import request, jsonify
from src.config import Config
from src.models.sparse_model import SparseModel

def register_health_routes(app):
    @app.route('/health', methods=['GET'])
    def health():
        return jsonify({
            'status': 'ok',
            'models': 'loaded',
            'batch_size': Config.batch_size,
            'max_batch_size': Config.max_batch_size
        })