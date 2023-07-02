import logging
from config import Config
from app.server import Server

def create_app():
    logging.info("Creating job-executor app")
    config = Config()
    s = Server()
    s.run(config)
    