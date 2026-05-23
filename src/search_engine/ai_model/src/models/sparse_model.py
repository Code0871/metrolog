from sklearn.feature_extraction.text import TfidfVectorizer
from pathlib import Path
from typing import List
import pickle
import os

class SparseModel:
    def __init__(self, model_name: str = "Qdrant/bm25", cache_dir: str = "./models_cache"):
        self.cache_dir = Path(cache_dir)
        self.vectorizer = TfidfVectorizer(
            analyzer='char_wb',
            ngram_range=(2, 5),
            max_features=100000,
            min_df=1,
            max_df=1.0,
            sublinear_tf=True
        )
        self.is_fitted = False
        self.bm25_index_path = os.getenv("bm25_index_path", "./bm25_index.pkl")
        
        # Загружаем индекс и обучаем модель
        self._load_and_fit_from_index()
    
    def _load_and_fit_from_index(self):
        """Загрузка индекса и обучение модели"""
        if os.path.exists(self.bm25_index_path):
            try:
                with open(self.bm25_index_path, 'rb') as f:
                    index = pickle.load(f)
                
                if 'corpus' in index:
                    documents = index['corpus']
                    print(f"Loaded {len(documents)} documents from corpus (old format)")
                elif 'documents' in index:
                    documents = index['documents']
                    print(f"Loaded {len(documents)} documents from index (new format)")
                else:
                    print(f"Unknown index format. Keys: {index.keys()}")
                    return
                
                if documents:
                    print(f"Training sparse model on {len(documents)} documents...")
                    self.vectorizer.fit(documents)
                    self.is_fitted = True
                    print(f"Sparse model trained! Vocabulary size: {len(self.vectorizer.vocabulary_)}")
                else:
                    print("Index loaded but no documents found")
            except Exception as e:
                print(f"Failed to load index: {e}")
        else:
            print(f"Index not found at {self.bm25_index_path}")
    
    def encode(self, texts: List[str]) -> List[dict]:
        """Получить sparse вектора для текстов"""
        if not self.is_fitted:
            print(f"Model not fitted, training on current texts...")
            self.vectorizer.fit(texts)
            self.is_fitted = True
        
        sparse_matrix = self.vectorizer.transform(texts)
        
        sparse_vectors = []
        for i in range(sparse_matrix.shape[0]):
            row = sparse_matrix[i]
            sparse_vectors.append({
                "indices": row.indices.tolist(),
                "values": row.data.tolist()
            })
        
        return sparse_vectors
    
    def update_index(self, new_corpus):
        """Добавление документов в индекс и переобучение модели"""
        # Загружаем существующий индекс
        if os.path.exists(self.bm25_index_path):
            with open(self.bm25_index_path, 'rb') as f:
                index = pickle.load(f)
            
            # Поддерживаем оба формата
            if 'corpus' in index:
                documents = index['corpus']
            elif 'documents' in index:
                documents = index['documents']
            else:
                documents = []
        else:
            index = {'documents': [], 'document_count': 0}
            documents = []
        
        if isinstance(new_corpus, str):
            new_corpus = [new_corpus]
        
        added_count = 0
        for doc in new_corpus:
            if doc not in documents:
                documents.append(doc)
                added_count += 1
        
        if added_count > 0:
            print(f"Retraining model on {len(documents)} documents...")
            self.vectorizer.fit(documents)
            self.is_fitted = True
            print(f"Model retrained! Vocabulary size: {len(self.vectorizer.vocabulary_)}")
        
        # Сохраняем в новом формате (конвертируем старый если нужно)
        index['documents'] = documents
        index['document_count'] = len(documents)
        
        with open(self.bm25_index_path, 'wb') as f:
            pickle.dump(index, f)
        
        return added_count, len(documents)


SparseEmbeddingModel = SparseModel