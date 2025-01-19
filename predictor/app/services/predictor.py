import os
import numpy as np
from keras.api.models import load_model
from sklearn.preprocessing import StandardScaler
import logging
class Predictor:
    def __init__(self, model_dir="app/services/models"):
        self.model_dir = model_dir

        self.stdsc = StandardScaler()


    def load_trained_model(self, sec_id):
        try:
            model_path = os.path.join(self.model_dir, f"{sec_id}.keras")
            model = load_model(model_path)
            return model
        except Exception as e:
            logging.error(f"Error loading model: {e}")
            return None

    def prepare_data_for_prediction(self, Data, train_len=120):
        try:
            data = Data.values
            data = data.reshape(data.shape[0], 1)
            std_data = self.stdsc.fit_transform(data)

            X_pred = []
            for i in range(train_len, std_data.shape[0]):
                X_pred.append(std_data[i - train_len: i, 0])
            X_pred = np.array(X_pred)
            X_pred = np.reshape(X_pred, (X_pred.shape[0], X_pred.shape[1], 1))

            return X_pred[-1].reshape(1, train_len, 1)
        except Exception as e:
            print(f"Error in prepare_data_for_prediction: {e}")
            return None

    def predict_growth(self, model, dataset, future_time=96):
        try:
            if model is None:
                raise ValueError("Model is not loaded.")

            if dataset is None or 'Цена закрытия' not in dataset.columns:
                raise ValueError("Dataset is required for prediction and must include 'Цена закрытия'.")

            X = self.prepare_data_for_prediction(dataset)
            logging.info(f"Shape of X_pred: {X.shape}")
            future_predictions = []
            for _ in range(future_time):
                prediction = model.predict(X)
                future_predictions.append(prediction[0, 0])
                new_prediction = np.array([[[prediction[0, 0]]]])
                X = np.append(X[:, 1:, :], new_prediction, axis=1)
            future_predictions = np.array(future_predictions).reshape(-1, 1)
            future_predictions = self.stdsc.inverse_transform(future_predictions)
            flat_array = np.array(future_predictions).flatten()
            flat_array = flat_array.tolist()
            return flat_array
        except Exception as e:
            logging.error(f"Error in predict_growth: {e}")
            return None






