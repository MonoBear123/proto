import logging
from datetime import datetime, timedelta
import  os
from proto_gen.predict_pb2_grpc import StonksPredictorServicer
from proto_gen.predict_pb2 import PredictorResponse
from app.services.db_manager import DBManager
from app.services.model_trainer import ModelTrainer
from app.services.predictor import Predictor

import traceback

TREANING_MODELS = []
MAX_AMOUNT_MODELS = 5

class StonksPredictorHandler(StonksPredictorServicer):
    def __init__(self, db_manager: DBManager):
        self.db_manager = db_manager
        self.model_trainer = ModelTrainer()
        self.predictor = Predictor()

    def Predictor(self, request, context):
        sec_id = request.query
        logging.info(f"Received secid: {sec_id}")
        if sec_id in TREANING_MODELS:
            return PredictorResponse(numbers=[1])
        try:
            week = self.model_trainer.get_dataset(sec_id,20)

            week_float = week["Закрытие"].astype(float).tolist()

            company_info  = self.db_manager.get_company_info(sec_id)
            if company_info is None:
                model_path = os.path.join(os.path.dirname(__file__), "../../services/models", f"{sec_id}.keras")
                model_path = os.path.normpath(model_path)
                if os.path.isfile(model_path):
                    model = self.predictor.load_trained_model(sec_id=sec_id)
                    data = self.model_trainer.get_dataset(sec_id)
                    predict = self.predictor.predict_growth(model,data)
                    self.db_manager.save_model_to_db(sec_id,model_path,predict)
                    res = week_float + predict
                    return PredictorResponse(numbers= res)

                logging.warning(f"secid '{sec_id}' not found in database.")
                logging.info(f"Fetching data for secid '{sec_id}'...")

                data = self.model_trainer.get_dataset(sec_id)
                if data is None:
                    logging.error("data in empty")
                    return PredictorResponse()

                logging.info("Training the model...")
                if len(TREANING_MODELS) == MAX_AMOUNT_MODELS:
                    return PredictorResponse(numbers=[0,0])
                TREANING_MODELS.append(sec_id)
                model, path = self.model_trainer.get_fit_model(data,sec_id=sec_id)
                if model is None:
                    TREANING_MODELS.remove(sec_id)
                    return PredictorResponse()
                logging.info("Model training completed.")
                prediction = self.predictor.predict_growth(model, data)
                self.db_manager.save_model_to_db(sec_id, path, prediction)
                TREANING_MODELS.remove(sec_id)
                res = week_float + prediction
                return PredictorResponse(numbers=res)

            else:
                id, last_prediction_time = company_info
                if last_prediction_time:

                    time_diff = datetime.now() - last_prediction_time

                    if time_diff > timedelta(hours=4):
                        logging.info(f"Last prediction was {time_diff} ago. Refreshing dataset and making new prediction.")
                        data = self.model_trainer.get_dataset(sec_id)
                        if data is None:
                            return PredictorResponse()
                        model = self.predictor.load_trained_model(sec_id=sec_id)
                        prediction = self.predictor.predict_growth(model, data)
                        self.db_manager.insert_prediction(sec_id, prediction)
                        res = week_float + prediction
                        return PredictorResponse(numbers= res)

                    else:
                        logging.info(f"Last prediction was {time_diff} ago. Using previous prediction.")
                        prediction = self.db_manager.get_predictions(sec_id)
                        res = week_float + prediction
                        return PredictorResponse(numbers= res)

        except Exception as _:
            logging.error(f"Error during prediction for secid '{sec_id}': {traceback.format_exc()}")
