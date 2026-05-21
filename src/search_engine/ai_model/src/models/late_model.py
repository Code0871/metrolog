from pathlib import Path
import torch
from typing import List
import numpy as np

class LateInteractionModel:
    def __init__(self, model_name: str, cache_dir: str = "./models_cache"):
        from transformers import AutoTokenizer, AutoModel
        
        self.model_name = model_name
        self.cache_dir = Path(cache_dir)
        
        model_path = self._get_model_path() or model_name
        
        print(f"Loading model from: {model_path}")
        
        self.tokenizer = AutoTokenizer.from_pretrained(
            str(model_path) if model_path else model_name,
            cache_dir=str(cache_dir)
        )
        self.model = AutoModel.from_pretrained(
            str(model_path) if model_path else model_name,
            cache_dir=str(cache_dir)
        )
        self.model.eval()
    
    def _get_model_path(self) -> Path:
        model_dir_name = "models--" + self.model_name.replace("/", "--")
        snapshots_dir = self.cache_dir / model_dir_name / "snapshots"
        
        if snapshots_dir.exists():
            snapshots = list(snapshots_dir.iterdir())
            if snapshots:
                return snapshots[0]
        return None
    
    def encode(self, texts: List[str]) -> List[torch.Tensor]:
        with torch.no_grad():
            inputs = self.tokenizer(
                texts,
                padding=True,
                truncation=True,
                max_length=512,
                return_tensors="pt"
            )
            
            outputs = self.model(**inputs)
            
            embeddings = outputs.last_hidden_state
            
            result = []
            for i, text in enumerate(texts):
                # Получаем длину без паддинга
                length = inputs['attention_mask'][i].sum()
                token_embeddings = embeddings[i, 1:length-1, :]  # убираем специальные токены
                result.append(token_embeddings)
            
            return result
    
    def maxsim_score(self, query_emb: torch.Tensor, doc_emb: torch.Tensor) -> float:

        similarity = torch.matmul(query_emb, doc_emb.T)
        # MaxSim: max по каждому токену запроса, затем сумма
        max_sim = similarity.max(dim=1).values.sum()
        return max_sim.item()
    
    def search(self, query: str, documents: List[str], top_k: int = 5) -> List[dict]:
        query_emb = self.encode([query])[0]
        doc_embs = self.encode(documents)
        
        scores = []
        for i, doc_emb in enumerate(doc_embs):
            score = self.maxsim_score(query_emb, doc_emb)
            scores.append({
                "index": i,
                "text": documents[i],
                "score": score
            })
        
        scores.sort(key=lambda x: x["score"], reverse=True)
        return scores[:top_k]