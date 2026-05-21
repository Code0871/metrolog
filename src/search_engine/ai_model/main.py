# main.py
from flask import Flask
from src.routes import dense, sparse, late, bm25, health

app = Flask(__name__)

# Регистрируем все routes
health.register_health_routes(app)
dense.register_dense_routes(app)
sparse.register_sparse_routes(app)
late.register_late_routes(app)
bm25.register_bm25_routes(app)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=False)