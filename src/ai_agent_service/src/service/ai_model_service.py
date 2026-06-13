import os
import catboost
import logging
from typing import List, Dict, Any, Union, Optional
import pickle
import pandas as pd
import psycopg2

class AIModel:
    def __init__(self, model_dir: Optional[str] = None):
        self.model_dir = model_dir or "ai_model"
        self.model = None
        self.load_model()
        self.features()
    
    def load_model(self):
        model = catboost.CatBoostClassifier()
        model.load_model(os.path.join(self.model_dir, "buyer_model.cbm"))
        self.model = model
        logging.info("AI model loaded successfully.")
    
    ## Удобная функция для отладки, помогает понять, какие признаки нужны модели
    def features(self):
        with open (os.path.join(self.model_dir, "model_files/model_features.pkl"), "rb") as f:
            features = pickle.load(f)
        
        print("Модель требует следующие признаки:", features)

    def predict(self, df: pd.DataFrame) -> Any:
        df_pred = df.copy()
        df_pred = df_pred.drop(columns=['miinstance_passport', 'miinstance_name', 'miinstance_type'])
        
        categorical_features = [
            'miinstance_state_condition',
            'miinstance_tech_condition', 
            'type_of_fact_mk',
            'type_of_planned_mk'
        ]
        numeric_features = [col for col in df_pred.columns if col not in categorical_features]
        
        for col in numeric_features:
            df_pred[col] = pd.to_numeric(df_pred[col], errors='coerce')
            df_pred[col] = df_pred[col].fillna(df_pred[col].mean())
        
        for col in categorical_features:
            if col in df_pred.columns:
                df_pred[col] = df_pred[col].fillna('unknown').astype(str)
        
        if self.model is None:
            raise ValueError("Model is not loaded.")
        
        predictions = self.model.predict(df_pred)
        
        if len(predictions.shape) > 1:
            predictions = predictions.flatten()
        
        df['prediction'] = predictions
        # print(df.head())
        return df
    
    ## TODO: переписать метод на использование http к основному сервису парка СИ
    ## Сейчас так из-за нехватки времени
    def get_data_and_transform_to_df(self, start_date: str, end_date: str) -> pd.DataFrame:
        conf = psycopg2.connect(
            dbname=os.getenv("main_base_dbname"),
            user=os.getenv("main_base_user"),
            password=os.getenv("main_base_password"),
            host=os.getenv("main_base_host"),
            port=os.getenv("main_base_port")
        )
        query = """
            select
                miinstance_passport,
                miinstance_name,
                miinstance_type,
                miinstance_state_condition,
                miinstance_tech_condition,
                type_of_fact_mk,
                type_of_planned_mk,
                ((now()::date - issue_date::date) / 365.25) as age_years,
                (now()::date - fact_date_of_mk::date) as days_since_last_mk,
                (plan_date_of_mk::date - now()::date) as days_until_planned_mk,
                case
                    when extract(day from (now()::date - fact_date_of_mk)) > (prmk * 30.44) then 1
                    else 0
                end as mpi_expired,
                case
                    when date_of_debit is not null then 1
                    else 0
                end as is_written_off	
            from miinstance
            where date_of_debit between %(start_date)s and %(end_date)s
                or (fact_date_of_mk + (prmk * 30.44 || ' days')::interval) between %(start_date)s and %(end_date)s
        """
        cursor = conf.cursor()
        cursor.execute(query, {"start_date": start_date, "end_date": end_date})
        data = cursor.fetchall()
        columns = [desc[0] for desc in cursor.description]
        df = pd.DataFrame(data, columns=columns)
        cursor.close()
        conf.close()
        return df
        

ai_model = AIModel()
df = ai_model.get_data_and_transform_to_df('2024-01-01','2026-12-31')
ai_model.predict(df)