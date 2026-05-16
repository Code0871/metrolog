from pathlib import Path
import os
from dotenv import load_dotenv
from huggingface_hub import snapshot_download
from sentence_transformers import SentenceTransformer
from optimum.onnxruntime import ORTModelForFeatureExtraction
from transformers import AutoTokenizer

path = Path(__file__).parents[1]

load_dotenv(dotenv_path=path / "config" / "config.env")
model_id = os.getenv("model_from_hugging_face")

model_dir = Path(__file__).parent / "model"

snapshot_download (
    repo_id = model_id,
    allow_patterns = ["*.json", "*.bin", "*.safetensors", "*.txt"],
    ignore_patterns=["*.msgpack", "*.h5"],
    local_dir = model_dir
)

onnx_path = Path(__file__).parent / "onnx_model"

# Загружаем PyTorch модель и конвертируем
model = ORTModelForFeatureExtraction.from_pretrained(
    model_id,
    export=True,
    provider="CPUExecutionProvider",  # Или "CUDAExecutionProvider" для GPU
)

tokenizer = AutoTokenizer.from_pretrained(model_id)

# Сохраняем ONNX модель и токенизатор
model.save_pretrained(onnx_path)
tokenizer.save_pretrained(onnx_path)