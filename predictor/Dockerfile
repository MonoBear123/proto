FROM python:3.9.13-slim

# Установка системных зависимостей
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev gcc && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копирование зависимостей
COPY requirements.txt .

# Обновление pip и установка зависимостей
RUN python -m pip install --upgrade pip && \
    pip cache purge && \
    pip install --no-cache-dir -r requirements.txt


# Копирование кода приложения
COPY . .

CMD ["python3", "main.py"]
