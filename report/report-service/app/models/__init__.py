import json

class JsonObject(object):
    def toJsonData(self):
        return self.toJsonStr().encode('utf-8')
    
    def toJsonStr(self):
        return json.dumps(self.to_dict())
    
    def to_dict(self):
        return vars(self)

from app.models._job import Job
from app.models._report import Report
from app.models._report import ReportList
from app.models._report import ReportDBInterface