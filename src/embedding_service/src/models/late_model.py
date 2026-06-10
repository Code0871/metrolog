from pathlib import Path
from pylate import models

class LateInteractionModel:
    def __init__(self, model_name: str = None, models_dir: str = None):
        if models_dir is None:
            self.cache_dir = Path(__file__).parents[2] / "models_dir" / "late"
        else:
            self.cache_dir = Path(models_dir)
        
        self.model_name = model_name or "answerai/answerai-colbert-small-v1"
        
        # Загружаем модель
        self.model = models.ColBERT(
            model_name_or_path=str(self.cache_dir),
            local_files_only=True  # Работаем только локально
        )
    
    def encode(self, texts, is_query=False):
        """Получить эмбеддинги для текстов"""
        return self.model.encode(texts, is_query=is_query)


late_model = LateInteractionModel()