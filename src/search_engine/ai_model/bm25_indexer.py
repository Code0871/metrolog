
## Временный файл для сохранения BM25 индекса (чисто для тестов)
import pickle
import psycopg2
from rank_bm25 import BM25Okapi
from datetime import datetime
import os

DB_CONFIG = {
    'host': 'localhost',
    'port': 5432,
    'database': 'miinstance_base',
    'user': 'postgres',
    'password': '2804'
}

def create_bm25_simple():
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()
        
        # Просто получаем все строки
        cursor.execute("SELECT miinstance_passport, miinstance_name FROM miinstance")
        rows = cursor.fetchall()
        
        # Получаем имена колонок
        col_names = [desc[0] for desc in cursor.description]
        print(f"Колонки: {col_names}")
        
        corpus = []
        doc_ids = []
        
        for row in rows:
            # Используем первую колонку как ID
            doc_ids.append(row[0])
            
            # Объединяем все строковые колонки
            text_parts = []
            for value in row:
                if value and isinstance(value, str):
                    text_parts.append(value)
            
            corpus.append(' '.join(text_parts))
        
        cursor.close()
        conn.close()
        
        if not corpus:
            print("Нет данных")
            return
        
        # Создаем индекс
        tokenized = [doc.lower().split() for doc in corpus]
        bm25 = BM25Okapi(tokenized)
        
        # Сохраняем
        with open('bm25_index.pkl', 'wb') as f:
            pickle.dump({'index': bm25, 'doc_ids': doc_ids, 'corpus': corpus}, f)
        
        print(f"✅ Индекс создан: {len(corpus)} документов")
        print(f"📁 Файл: {os.path.abspath('bm25_index.pkl')}")
        
    except Exception as e:
        print(f"Ошибка: {e}")

if __name__ == "__main__":
    create_bm25_simple()