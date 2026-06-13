import psycopg2
from psycopg2 import OperationalError
from config.config import db_host, db_user, db_password, db_name, db_port

class Database:
    def __init__(self, dbname=None, user=None, password=None, host=None, port=None):
        # Берем значения из config, если не переданы явно
        self.dbname = dbname or db_name
        self.user = user or db_user
        self.password = password or db_password
        self.host = host or db_host
        self.port = port or db_port
        
        try:
            self.conn = psycopg2.connect(
                dbname=self.dbname,
                user=self.user,
                password=self.password,
                host=self.host,
                port=self.port
            )
            print(f"Connection to database {self.dbname} successful!")
        except OperationalError as e:
            print(f"Connection to database failed! Error: {e}")
            self.conn = None
    
    def close(self):
        if self.conn:
            self.conn.close()
            print("Database connection closed.")
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()

# Автоматически берет параметры из config
db = Database()