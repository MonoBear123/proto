import logging

def setup_logger(config):
    logging.basicConfig(
        level=config.get('level', 'INFO'),
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )