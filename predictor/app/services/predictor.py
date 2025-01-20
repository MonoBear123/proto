import os
import numpy as np
import logging
from keras.api.models import load_model
from sklearn.preprocessing import MinMaxScaler

class Predictor:
    def __init__(self, model_dir="app/services/models"):
        self.model_dir = model_dir
        self.scaler = MinMaxScaler(feature_range=(0,1))
        self.sequence_size = 60
    def load_trained_model(self, sec_id):
        try:
            model_path = os.path.join(self.model_dir, f"{sec_id}.keras")
            model = load_model(model_path)
            return model
        except Exception as e:
            logging.error(f"Error loading model: {e}")
            return None

    def prepare_data_for_prediction(self, data, time_steps = 60):
        try:
            close_prices = data[['Закрытие']].values
            close_prices_scaled = self.scaler.fit_transform(close_prices)
            last_sequence = close_prices_scaled[-time_steps:]
            input_data = np.reshape(last_sequence, (1, time_steps, 1))
            return input_data
        except Exception as e:
            print(f"Error in prepare_data_for_prediction: {e}")
            return None

    def predict_growth(self, model, dataset):
        try:
            if model is None:
                raise ValueError("Model is not loaded.")

            if dataset is None or 'Закрытие' not in dataset.columns:
                raise ValueError("Dataset is required for prediction and must include 'Цена закрытия'.")

            input_data = self.prepare_data_for_prediction(dataset)
            predicted_scaled = model.predict(input_data)
            predicted_prices = self.scaler.inverse_transform(predicted_scaled)

            pred = [float(price) for price in predicted_prices[0]]
            return pred

        except Exception as e:
            logging.error(f"Error in predict_growth: {e}")
            return None