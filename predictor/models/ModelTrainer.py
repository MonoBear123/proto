import requests
from datetime import datetime, timedelta
import numpy as np
import pandas as pd
from sklearn.preprocessing import StandardScaler

from keras.api.models import Sequential
from keras.api.layers import Dense, LSTM

class ModelTrainer():
    def get_dataset(name_on_moscow_exchange: str, days_in_data: int = 90):
        end_date = datetime.now()
        start_date = end_date - timedelta(days=days_in_data)

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
            print(data)
            if 'candles' in data and 'data' in data['candles']:
                rows = data['candles']['data']
                all_data.extend(rows)
                if len(rows) < limit:
                    break
                params['start'] += limit
            else:
                break

        columns = data['candles']['columns']
        df = pd.DataFrame(all_data, columns=columns)
        df_filtered = df[['begin', 'close']].copy()
        df_filtered.columns = ['Дата', 'Цена закрытия']
        df_filtered['Дата'] = pd.to_datetime(df_filtered['Дата'])
        df_filtered = df_filtered.sort_values('Дата')
        df_filtered.to_csv(f"content/{name_on_moscow_exchange}_hourly_quotes.csv", index=False, encoding='utf-8-sig')
        return pd.read_csv(f'content/{name_on_moscow_exchange}_hourly_quotes.csv', index_col='Дата',
                           parse_dates=['Дата'])

    def prepairing_dataset_to_train(self, Data):
        train_data = Data.iloc[0:-1].values
        train_data = train_data.reshape(train_data.shape[0], 1)
        std_train_data = StandardScaler().fit_transform(train_data)
        X_train = []
        Y_train = []
        train_len = 120
        for i in range(train_len, std_train_data.shape[0]):
            X_train.append(std_train_data[i - train_len: i, 0])
            Y_train.append(std_train_data[i, 0])
        X_train, Y_train = np.array(X_train), np.array(Y_train)
        X_train = np.reshape(X_train, (X_train.shape[0], X_train.shape[1], 1))
        return X_train, Y_train


    def get_fit_model(self, Data, model=None, epochs: int = 15, batch_size=1, optimizer='Adam', loss='mean_squared_error'):
        X_train, Y_train = self.prepairing_dataset_to_train(Data)

        if not model:
            model = Sequential()
            model.add(LSTM(units=50, return_sequences=True, input_shape=(X_train.shape[1], 1)))
            model.add(LSTM(units=50, return_sequences=True))
            model.add(LSTM(units=50, return_sequences=True))
            model.add(LSTM(units=50))
            model.add(Dense(units=1))

        model.compile(optimizer=optimizer, loss=loss)

        model.fit(X_train, Y_train, epochs=epochs, batch_size=batch_size)

        return model


data = ModelTrainer.get_dataset("SBER")