from concurrent import futures
import logging
import grpc
from app.services.db_manager import DBManager
from proto_gen import predict_pb2_grpc
from app.api.grpc.handlers import StonksPredictorHandler
from app.utils.config_loader import Config

from pathlib import Path

SERVICE_PATH = Path(__file__).parent


def main() -> None:
    config_path = SERVICE_PATH / "config/config.yaml"
    config_loader = Config(str(config_path))

    try:
        db_name = config_loader.get('database.name')
        db_user = config_loader.get('database.user')
        db_password = config_loader.get('database.password')
        db_host = config_loader.get('database.host')
        db_port = config_loader.get('database.port')

        server_port = config_loader.get('server.port', 42020)

        db = DBManager(db_name, db_user, db_password, db_host, db_port)

        logging.info("Connected to the database.")
        db.create_tables()
        logging.info("Tables were been created.")

    except (FileNotFoundError, RuntimeError) as e:
        logging.error("Error during initialization: %s", e)
        return

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    handler = StonksPredictorHandler(db)

    predict_pb2_grpc.add_StonksPredictorServicer_to_server(handler, server)

    listen_addr = f"[::]:{server_port}"
    server.add_insecure_port(listen_addr)

    logging.info("Starting server on %s", listen_addr)

    server.start()  # Use `await` to ensure async server start
    server.wait_for_termination()  # Wait for server termination async

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    main()
