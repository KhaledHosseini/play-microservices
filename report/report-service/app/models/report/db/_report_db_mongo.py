import logging

from app.database_engines import MongoDB
from app.database_engines.mongo import MongoPagination
from app.models import Report
from app.models import ReportList
from app.models import ReportDB
from config import Config


class ReportDBMongo(MongoDB, ReportDB):
    def __init__(self, cfg: Config):
        super().__init__(cfg)
        ReportDB.__init__(self)
        self.dbName = cfg.DatabaseDBName
        self.collection = "ReportsCollection"
        self.db = self.client[self.dbName]

    def create(self, report: Report) -> Report:
        logging.info("creating report in databse.")
        try:
            self.db[self.collection].insert_one(report.to_dict())
        except Exception as e:
            logging.error("error creatin report in database: ", e)
            raise e

    def list(self, type: int, page: int, size: int) -> ReportList:
        logging.info("getting reports list...")
        try:
            find = {"type": type}
            coll = self.db[self.collection]
            doc_count = coll.count_documents(find)
            if doc_count == 0:
                return ReportList()
            pagination = MongoPagination(size=size, page=page)
            limit = pagination.getLimit()
            skip = pagination.getOffset()
            documents = coll.find().limit(limit).skip(skip)
            reports = [Report.from_dict(doc) for doc in documents]
            logging.info(f"results found: {reports}")
            return ReportList(
                totalCount=doc_count,
                totalPages=int(pagination.getTotalPages(int(doc_count))),
                page=int(pagination.getPage()),
                size=int(pagination.getSize()),
                hasMore=pagination.getHasMore(int(doc_count)),
                reports=reports,
            )
        except Exception as e:
            raise e