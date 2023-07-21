from config import Config
from app.email import EmailSender
from app.models.email_job.message_broker import MessageBrokerService

class Server:
    def run(self, cfg: Config):
        es = EmailSender(cfg)
        ecs = MessageBrokerService(cfg,es)
        ecs.run()
