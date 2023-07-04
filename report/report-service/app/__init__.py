from config import Config
from app.server import Server

import logging
def create_app():
    logging.info("creating app")
    config = Config()
    s = Server()
    s.run(config)
    