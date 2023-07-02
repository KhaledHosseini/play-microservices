import json
import logging
# from abc import ABC, abstractmethod

class JsonObject:
    def toJsonData(self):
        return json.dumps(vars(self)).encode('utf-8')
    def valueForKey(self,key,json):
        try:
            return json[key]
        except Exception as e:
            logging.error(f"no value found for json key : {key}")
            raise e

class EmailJob(JsonObject):
    def __init__(self, json):
        logging.info("Initilizing EmailJob from json object: ", json)
        self.jobId = self.valueForKey("jobId",json)
        self.status = int(self.valueForKey("status",json))
        self.JobData = self.valueForKey("jobData",json)
    
    def get_email(self):
        return Email(self.JobData)


class Email(JsonObject):
    def __init__(self,json):
        self.SourceAddress = self.valueForKey("SourceAddress",json)
        self.DestinationAddress= self.valueForKey("DestinationAddress",json)
        self.Subject = self.valueForKey("Subject",json)
        self.Message = self.valueForKey("Message",json)