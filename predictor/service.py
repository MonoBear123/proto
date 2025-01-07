from concurrent import futures
import logging
import grpc
import predict_pb2
import predict_pb2_grpc

# Определение класса, который будет обрабатывать gRPC-запросы
class StonksPredictorServicer(predict_pb2_grpc.StonksPredictorServicer):
    # Определение метода для обработки gRPC-запросов
    def Predictor(self, request, context) -> predict_pb2.PredictorResponse():
        # Вот тут нужен вызов функции функции,которая возвращает предсказание в виде листа с данными
        return predict_pb2.PredictorResponse(numbers=[12.0, 12.3])  # Пример

# Функция для запуска gRPC-сервера
def serve() -> None:
    # Создание и настройка grpc сервера
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))

    # Регистрируем сервис, добавляя обработчик для gRPC
    predict_pb2_grpc.add_StonksPredictorServicer_to_server(StonksPredictorServicer(), server)

    listen_addr = "[::]:42020"
    server.add_insecure_port(listen_addr)

    logging.info("Starting server on %s", listen_addr)

    # Запуск сервера
    server.start()

    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    serve()
