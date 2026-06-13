import os
from dotenv import load_dotenv

class Config:
    def __init__(self):
        load_dotenv()
        
        self.db_host = os.getenv('db_host', 'localhost')
        self.database_user = os.getenv("database_user", "postgres")
        self.database_name = os.getenv("database_name", "procurement_plan")
        self.database_password = os.getenv("database_password")
        self.database_url = os.getenv(
            "database_url",
            f"postgresql://{self.database_user}:{self.database_password}@{self.db_host}:5432/{self.database_name}"
        )
        self.main_base_host = os.getenv("main_base_host", "localhost")
        self.main_base_port = os.getenv("main_base_port", "5432")
        self.main_base_user = os.getenv("main_base_user", "postgres")
        self.main_base_password = os.getenv("main_base_password")
        self.main_base_dbname = os.getenv("main_base_dbname", "miinstance_base")

# Использование
config = Config()
print(config.database_url)