import os
import numpy as np
from keras.api.models import load_model
from sklearn.preprocessing import StandardScaler
from datetime import datetime

class Predictor:
    def __init__(self, model_dir='models'):
        self.model_dir = model_dir

    def predict_growth(self, model=None, model_path=None, dataset=None) -> [float]:
        try:
            if model is None:
                if model_path is None:
                    # Load the most recent model from the directory
                    models = [f for f in os.listdir(self.model_dir) if f.endswith('.h5')]
                    if not models:
                        raise FileNotFoundError("No model found.")
                    model_path = sorted(models, key=lambda x: os.path.getmtime(os.path.join(self.model_dir, x)))[-1]
                    model_path = os.path.join(self.model_dir, model_path)
                model = load_model(model_path)

            if dataset is None:
                raise ValueError("Dataset is required for prediction.")

            # Preprocess the dataset (reshape and scale as needed)
            data = dataset.values.reshape(-1, 1)
            std_data = StandardScaler().fit_transform(data)
            X = []
            Y = []
            train_len = 120
            for i in range(train_len, std_data.shape[0]):
                X.append(std_data[i - train_len: i, 0])
                Y.append(std_data[i, 0])
            X, Y = np.array(X), np.array(Y)
            X = np.reshape(X, (X.shape[0], X.shape[1], 1))

            # Predict on future prices
            future_predictions = model.predict(X)

            # Calculate growth (difference between consecutive predictions)
            growth = future_predictions[1:] - future_predictions[:-1]

            return growth
        except Exception as e:
            print(f"Error in predict_growth: {e}")
            return None
