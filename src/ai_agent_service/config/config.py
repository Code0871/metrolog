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

# Использование
config = Config()
print(config.database_url)