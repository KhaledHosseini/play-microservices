from smtplib import SMTP
from app.models.email_job import Email
from config import Config
import logging
import threading

class EmailSender:
    def __init__(self,config: Config):
        self.cfg = config
        self.lock = threading.Lock()
    
    def send(self, email: Email):

        logging.info("Sending email...")
        
        self.lock.acquire()

        sender = self.cfg.SmtpUser + "@" + self.cfg.EMailDomain
        message_template = ("""From: <{0}>
        To: <{1}>
        Subject: {2} 
        {3}""")
        host = self.cfg.MailServerHost
        try:
            with SMTP(host=host,port=25) as smtp:
                smtp.login(sender, self.cfg.SmtpPassord)
                smtp.sendmail(from_addr=sender,to_addrs= [email.DestinationAddress], msg= message_template.format(sender, email.DestinationAddress, email.Subject,email.Message))
        except Exception as e:
            logging.error("Sending email failed with error:", str(e))
            raise e
        finally:
            self.lock.release()