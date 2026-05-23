
## Временный файл для сохранения BM25 индекса
import pickle
from rank_bm25 import BM25Okapi
import os


def create_bm25_simple():
    if os.path.exists('bm25_index.pkl'):
        print("Индекс BM25 существует")
        return
    
    try:        
        corpus = ["вольтометр", "Амперметр", "Адгезиметр", "датчик", "анализатор", 
                  "измеритель", "Ампервольтваттметр", "Аспиратор", 
                  "Блок", "pH", "Вакуумметр", "Варметр", "Ваттметр", "Весы", "Водосчетчик", 
                  "Вольтметр", "Газоанализатор", "Газосигнализатор", "Измеритель", "Индикатор", 
                  "Калибратор", "Клещи", "Манометр", "Тахометр", "Мультиметр", "Прибор", "Термометр", 
                  "Тестер", "Система", "секундомер", "регулятор", "Регистратор", "Расходомер", "Приспособление", 
                  "пирометр", "Осциллограф", "Омметр", "Нивелир", "Оксиметр"]
        
        if not corpus:
            print("Нет данных")
            return
        
        ## Создаем индекс
        tokenized = [doc.lower().split() for doc in corpus]
        bm25 = BM25Okapi(tokenized)
        
        ## Сохраняем
        with open('bm25_index.pkl', 'wb') as f:
            pickle.dump({'index': bm25, 'corpus': corpus}, f)
        
        print(f"Индекс создан: {len(corpus)} документов")
        print(f"Файл: {os.path.abspath('bm25_index.pkl')}")
        
    except Exception as e:
        print(f"Ошибка: {e}")

def get_bm25_indexer():
    """Загружает BM25 индекс из файла"""
    if os.path.exists('bm25_index.pkl'):
        with open('bm25_index.pkl', 'rb') as f:
            data = pickle.load(f)
            return data['index']
    else:
        # Если индекса нет, создаем
        create_bm25_simple()
        with open('bm25_index.pkl', 'rb') as f:
            data = pickle.load(f)
            return data['index']


bm25_indexer = get_bm25_indexer()