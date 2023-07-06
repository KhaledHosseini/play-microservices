import logging
from app.architecture import Context, Containers,Contract
def create_app():
    logging.info("creating app...")
    Context().create()
    Containers().create()