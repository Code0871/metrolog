from flask import request, jsonify
from src.models.sparse_model import SparseModel

def register_bm25_routes(app):
    @app.route('/update_index', methods=["POST"])
    def update_index():
        try:
            data = request.json
            new_corpus = data.get('miinstance_name')
            
            if not new_corpus:
                return jsonify({"error": "Нет данных для добавления"}), 400
            
            sparse_model = SparseModel()
            added_count, total_count = sparse_model.update_index(new_corpus)
            
            if added_count == total_count:
                return jsonify({
                    "message": "Индекс создан",
                    "documents_count": total_count
                }), 200
            else:
                return jsonify({
                    "message": "Индекс обновлён",
                    "previous_count": total_count - added_count,
                    "added_count": added_count,
                    "total_count": total_count
                }), 200
                
        except Exception as e:
            return jsonify({"error": str(e)}), 500