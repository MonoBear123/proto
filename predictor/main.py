from concurrent import futures
import logging
import grpc
from grpc import predict_pb2
from grpc import predict_pb2_grpc
from models.ModelTrainer import ModelTrainer

class StonksPredictorServicer(predict_pb2_grpc.StonksPredictorServicer):
    def Predictor(self, request, context) -> predict_pb2.PredictorResponse():
        # Вот тут нужен вызов функции,которая возвращает предсказание в виде листа с данными

        return predict_pb2.PredictorResponse(numbers=[12.0, 12.3])  # Пример


def main() -> None:
    # Создание и настройка grpc сервера
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))

    # Регистрирация сервера, добавляя обработчик для gRPC
    predict_pb2_grpc.add_StonksPredictorServicer_to_server(StonksPredictorServicer(), server)

    listen_addr = "[::]:42020"
    server.add_insecure_port(listen_addr)

    logging.info("Starting server on %s", listen_addr)

    server.start()

    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    main()
