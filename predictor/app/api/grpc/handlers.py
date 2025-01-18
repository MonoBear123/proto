import logging
from datetime import datetime, timedelta

from proto_gen.predict_pb2_grpc import StonksPredictorServicer
from proto_gen.predict_pb2 import PredictorResponse
from app.services.db_manager import DBManager
from app.services.model_trainer import ModelTrainer
from app.services.predictor import Predictor
TREANING_MODELS = []
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
            company_info = self.db_manager.get_company_info(sec_id)
            if company_info is None:
                logging.warning(f"secid '{sec_id}' not found in database.")
                logging.info(f"Fetching data for secid '{sec_id}'...")

                data = self.model_trainer.get_dataset(sec_id)
                if data is None:
                    logging.error("data in empty")
                    return PredictorResponse()

                logging.info("Training the model...")
                TREANING_MODELS.append(sec_id)
                model, path = self.model_trainer.get_fit_model(data,sec_id=sec_id)
                if model is None:
                    TREANING_MODELS.remove(sec_id)
                    return PredictorResponse()
                logging.info("Model training completed.")
                prediction = self.predictor.predict_growth(model, data)
                logging.info(prediction)
                self.db_manager.save_model_to_db(sec_id, path, prediction)
                TREANING_MODELS.remove(sec_id)
                return PredictorResponse(numbers=prediction)

            else:
                id, last_prediction_time = company_info
                if last_prediction_time:
                    last_prediction_time_obj = datetime.fromisoformat(last_prediction_time)
                    time_diff = datetime.now() - last_prediction_time_obj

                    if time_diff > timedelta(hours=4):
                        logging.info(f"Last prediction was {time_diff} ago. Refreshing dataset and making new prediction.")
                        data = self.model_trainer.get_dataset(sec_id)
                        if data is None:
                            return PredictorResponse()

                        model_path = self.db_manager.get_model_path(sec_id)
                        prediction = self.predictor.predict_growth(model_path, data)
                        self.db_manager.insert_prediction(sec_id, prediction)
                        return PredictorResponse(numbers=prediction)

                    else:
                        logging.info(f"Last prediction was {time_diff} ago. Using previous prediction.")
                        prediction = self.db_manager.get_predictions(id)
                        return PredictorResponse(numbers=prediction)

        except Exception as e:
            logging.error(f"Error during prediction for secid '{sec_id}': {e}")
            return PredictorResponse()  # return empty result in case of error
