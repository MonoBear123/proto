import logging
from sklearn.preprocessing import MinMaxScaler
import requests
from datetime import datetime, timedelta
from sklearn.preprocessing import StandardScaler
import requests
from datetime import datetime, timedelta
import math
import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
from sklearn.model_selection import train_test_split
from keras.api.callbacks import ModelCheckpoint
from keras.api.models import Sequential
from keras.api.layers import Dense, LSTM, Dropout
import os



class ModelTrainer():
    def __init__(self):
        self.scaler =  MinMaxScaler(feature_range=(0, 1))

    @staticmethod
    def get_dataset(name_on_moscow_exchange: str, days_in_data: int = 60000):
        try:
            end_date = datetime.now()
            start_date = end_date - timedelta(days = days_in_data)

            start_date_str = start_date.strftime('%Y-%m-%d')
            end_date_str = end_date.strftime('%Y-%m-%d')

            url = f"https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities/{name_on_moscow_exchange}/candles.json"

            params = {
                'from': start_date_str,
                'till': end_date_str,
                'interval': 60,
                'start': 0
            }
            all_data = []
            limit = 500

            while True:
                response = requests.get(url, params=params)
                data = response.json()
                if 'candles' in data and 'data' in data['candles']:
                    rows = data['candles']['data']
                    all_data.extend(rows)
                    if len(rows) < limit:
                        break
                    params['start'] += limit
                else:
                    break
            if not os.path.exists("/content"):
                os.makedirs("/content")
            columns = data['candles']['columns']
            df = pd.DataFrame(all_data, columns=columns)
            print(df.columns)
            df_filtered = df[['begin','open', 'high', 'low', 'close', 'volume']].copy()
            df_filtered.columns = ['Дата','Цена открытия', 'Пик', 'Минимум', 'Закрытие', 'volume']
            df_filtered['Дата'] = pd.to_datetime(df_filtered['Дата'])
            df_filtered = df_filtered.sort_values('Дата')
            df_filtered.to_csv(f"/content/{name_on_moscow_exchange}_hourly_quotes.csv", index=False, encoding='utf-8-sig')
            return pd.read_csv(f'/content/{name_on_moscow_exchange}_hourly_quotes.csv', parse_dates=['Дата'])
        except Exception as e:
            logging.error(f"Ошибка в функции get_dataset: {e}")
            return None
    def prepare_data_to_fit(self,data):
        try:
            data = data.dropna()
            close_price = data[['Закрытие']].values
            close_price_scaled = self.scaler.fit_transform(close_price)

            time_steps = 60
            future_steps = 24

            X = []
            y = []

            for i in range(time_steps, len(close_price_scaled) - future_steps + 1):
                X.append(close_price_scaled[i - time_steps:i, 0])
                y.append(close_price_scaled[i:i + future_steps, 0])

            X, y = np.array(X), np.array(y)

            X = np.reshape(X, (X.shape[0], X.shape[1], 1))

            X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, shuffle=False)
            return X_train, X_test,y_train,y_test
        except Exception as e:
            logging.error(f"Error in prepairing_dataset_to_train: {e}")
            return None, None


    def get_fit_model(self,data,future_steps = 24,sec_id = None):
        try:
            model_dir = "app/services/models"
            if not os.path.exists(model_dir):
                os.makedirs(model_dir, exist_ok=True)

            model_save_path = os.path.join(model_dir, f"{sec_id}.keras")
            best_model_checkpoint_callback = ModelCheckpoint(model_save_path,monitor="val_loss",save_best_only=True,mode="min",verbose=0)
            X_train, X_test, y_train, y_test = ModelTrainer.prepare_data_to_fit(data)
            regressor = Sequential()
            regressor.add(LSTM(units=100, return_sequences=True, input_shape=(X_train.shape[1], 1)))
            regressor.add(Dropout(0.2))
            regressor.add(LSTM(units=100, return_sequences=True))
            regressor.add(Dropout(0.2))
            regressor.add(LSTM(units=100, return_sequences=True))
            regressor.add(Dropout(0.2))
            regressor.add(LSTM(units=100))
            regressor.add(Dropout(0.2))
            regressor.add(Dense(units=future_steps))
            regressor.compile(optimizer='adam', loss='mean_squared_error')
            regressor.fit(X_train, y_train, epochs=200, batch_size=32, validation_data=(X_test, y_test), callbacks = [best_model_checkpoint_callback])
            return regressor
        except Exception as e:
            logging.error(f"Error in get_fit_model: {e}")
            return None, None

