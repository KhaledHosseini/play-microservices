from config import Config
from app.email.email_sender import EmailSender
from app.models.email_job.message_broker.emailjob_consumer_service import EmailjobConsumerService

class Server:
    def run(self, cfg: Config):
        es = EmailSender(cfg)
        ecs = EmailjobConsumerService(cfg,es)
        ecs.run()
