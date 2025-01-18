import psycopg2
from datetime import datetime

class DBManager:
    def __init__(self, db_name, user, password, host, port):
        self.conn = psycopg2.connect(
            dbname=db_name,
            user=user,
            password=password,
            host=host,
            port=port
        )

    def create_tables(self):
        cursor = self.conn.cursor()
        cursor.execute("""
        CREATE TABLE IF NOT EXISTS predictions (
            id SERIAL PRIMARY KEY,
            company_id TEXT,
            prediction REAL[] NOT NULL,
            last_training_time TIMESTAMP NOT NULL,
            last_prediction_time TIMESTAMP NOT NULL,
            model_path TEXT
        );
        """)
        self.conn.commit()

    def get_company_info(self, sec_id):
        try:
            cursor = self.conn.cursor()
            cursor.execute("SELECT * FROM predictions WHERE sec_id = %s;", (sec_id,))
            result = cursor.fetchone()
            return result
        except Exception:
            return None

    def insert_prediction(self, company_id, prediction):
        cursor = self.conn.cursor()
        last_prediction_time = datetime.now()
        cursor.execute("""
        INSERT INTO predictions (company_id, prediction, last_prediction_time)
        VALUES (%s, %s, %s)
        ON CONFLICT (company_id)
        DO UPDATE SET prediction = EXCLUDED.prediction, last_prediction_time = EXCLUDED.last_prediction_time;
        """, (company_id, prediction, last_prediction_time))
        self.conn.commit()

    def get_model_path(self, company_id: int) -> str:
        cursor = self.conn.cursor()
        cursor.execute("SELECT model_path FROM predictions WHERE company_id = %s;", (company_id,))
        result = cursor.fetchone()
        return result[0] if result else None

    def get_predictions(self, company_id) -> [float]:
        cursor = self.conn.cursor()
        cursor.execute("SELECT prediction FROM predictions WHERE company_id = %s;", (company_id,))
        results = cursor.fetchall()
        return [float(row[0]) for row in results]

    def save_model_to_db(self, company_id: int, model_path: str, prediction: list):
        cursor = self.conn.cursor()
        last_prediction = datetime.now()
        last_trained = datetime.now()
        cursor.execute("""
        INSERT INTO predictions (company_id, model_path, prediction, last_training_time, last_prediction_time)
        VALUES (%s, %s, %s, %s, %s)
        ON CONFLICT (company_id)
        DO UPDATE SET
            model_path = EXCLUDED.model_path,
            prediction = EXCLUDED.prediction,
            last_training_time = EXCLUDED.last_training_time,
            last_prediction_time = EXCLUDED.last_prediction_time;
        """, (company_id, model_path, prediction, last_trained, last_prediction))
        self.conn.commit()

