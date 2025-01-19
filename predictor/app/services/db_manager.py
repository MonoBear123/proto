import logging
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
            company_id TEXT UNIQUE,
            prediction REAL[] NOT NULL,
            last_training_time TIMESTAMP NOT NULL,
            last_prediction_time TIMESTAMP NOT NULL,
            model_path TEXT
        );
        """)
        self.conn.commit()

    def get_company_info(self, sec_id: str):
        try:
            cursor = self.conn.cursor()
            cursor.execute("SELECT * FROM predictions WHERE company_id = %s;", (sec_id,))
            result = cursor.fetchone()
            return result[0], result[3]
        except Exception:
            return None

    def insert_prediction(self, company_id: str, prediction):
        cursor = self.conn.cursor()
        last_prediction_time = datetime.now()
        cursor.execute("""
            UPDATE predictions
            SET  prediction = %s,  last_prediction_time = %s
            WHERE company_id = %s;
            """, ( prediction,  last_prediction_time, company_id))
        logging.info(f"Updated record for company_id: {company_id}")
        self.conn.commit()

    def get_model_path(self, company_id: str) -> str:
        cursor = self.conn.cursor()
        cursor.execute("SELECT model_path FROM predictions WHERE company_id = %s;", (company_id,))
        result = cursor.fetchone()
        return result[0] if result else None

    def get_predictions(self, company_id) -> [float]:
        cursor = self.conn.cursor()
        cursor.execute("SELECT prediction FROM predictions WHERE company_id = %s;", (company_id,))
        results = cursor.fetchall()
        logging.info(results)
        logging.info(results[0][0])

        res = [float(row) for row in results[0][0]]
        logging.info(res)
        return res

    def save_model_to_db(self, company_id: str, model_path: str, prediction: list):
        cursor = self.conn.cursor()
        last_prediction = datetime.now()
        last_trained = datetime.now()
        try:
            cursor.execute("""
        INSERT INTO predictions (company_id, model_path, prediction, last_training_time, last_prediction_time)
        VALUES (%s, %s, %s, %s, %s)
        """, (company_id, model_path, prediction, last_trained, last_prediction))
            self.conn.commit()
            logging.info("save model",company_id,model_path,last_prediction)
        except Exception as e:
            logging.error(f"Error during database operation: {e}")
            self.conn.rollback()

