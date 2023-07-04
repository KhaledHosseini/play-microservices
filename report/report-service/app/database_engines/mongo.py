from pymongo import MongoClient
import math
from config import Config

class MongoDB:
    def __init__(self, cfg: Config) -> None:
        host = cfg.DatabaseHostName
        port = cfg.DatabasePort
        user = cfg.DatabaseUser
        password = cfg.DatabasePass
        schema = cfg.DatabaseSchema
        connection_string = f"{schema}://{user}:{password}@{host}:{port}"
        # Create a MongoClient object
        self.client = MongoClient(connection_string)

class MongoPagination:
    def __init__(self, size: int, page:int) -> None:
        self.size = max(size,1)
        self.page = page
    def getOffset(self):
        if self.page == 0:
            return 0
        return (self.page - 1) * self.size
    def getPage(self):
        return self.page
    def getLimit(self):
        return self.size
    def getSize(self):
        return self.size
    def getTotalPages(self, totalCount: int):
        d = float(totalCount) / float(self.getSize())
        return int(math.ceil(d))
    def getHasMore(self,totalCount: int):
        return self.getPage() < totalCount/self.getSize()
