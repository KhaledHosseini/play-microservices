import logging

from app.database_engines import MongoDB
from app.models import Report
from app.models import ReportList
from app.models import ReportDBInterface
from config import Config


class ReportDBMongo(MongoDB, ReportDBInterface):
    def __init__(self, cfg: Config):
        super().__init__(cfg)
        ReportDBInterface.__init__(self)
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

    def list(self, page: int, size: int) -> ReportList:
        logging.info("getting reports list...")
        try:
            coll = self.db[self.collection]
            doc_count = coll.count_documents({})
            if doc_count == 0:
                logging.info(f"list count is zero...")
                return ReportList()
            pagination = MongoDB.MongoPagination(size=size, page=page)
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
            logging.info(f"error getting reports list...{e}")
            raise e