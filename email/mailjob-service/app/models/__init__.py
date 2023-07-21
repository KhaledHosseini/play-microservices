import json

class JsonObject:
    def toJsonData(self):
        return self.toJsonStr().encode('utf-8')
    
    def toJsonStr(self):
        return json.dumps(self.to_dict())
    
    def to_dict(self):
        return vars(self)

from app.models._email_job import EmailJob
from app.models._email_job import Email
from app.models._email_job import JobStatusEnum