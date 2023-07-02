from config import Config
from app.models.message_broker import run_job_consumers

def create_app():
    config = Config()
    run_job_consumers(config)
    