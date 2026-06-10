import logging
from flask import Flask
from src.config import Config
from src.routes.dense import register_dense_routes
from src.routes.sparse import register_sparse_routes
from src.routes.late import register_late_routes
from src.routes.health import register_health_routes
from src.routes.bm25 import register_bm25_routes

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def create_app():
    app = Flask(__name__)
    
    # Регистрируем все роуты
    register_dense_routes(app)
    register_sparse_routes(app)
    register_late_routes(app)
    register_health_routes(app)
    register_bm25_routes(app)
    
    return app