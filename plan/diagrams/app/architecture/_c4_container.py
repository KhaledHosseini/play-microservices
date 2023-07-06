from diagrams import Cluster, Diagram
from diagrams.onprem.queue import Kafka
from diagrams.onprem.inmemory import Redis
from diagrams.onprem.database import Mongodb, Postgresql
from diagrams.c4 import Person, Container, Database, System, SystemBoundary, Relationship
class Containers:
    def create(self):
        with Diagram("Containers", show=False, direction="LR"):
            with Cluster(""):
                person = Person(name='users/admins', description='Any one whos has the application.',external=True)
                emaiSystem = System(name='emailsystem',description='Runs email-sending jobs')
                with SystemBoundary("Scheduler microservice app"):
                    client = Container(
                    name="Client Service",
                    technology="Typescript + React",
                    description="Application ui that users interact with",)

                    api_gateway = Container(
                    name="API-Gateway Service",
                    technology="Go + Rest",
                    description="Entrypoint of microservice application",)

                    auth = Container(
                    name="Authentication Service",
                    technology="Rust + gRPC",
                    description="Users Service responsible for creating users and authenticate them",)
                    auth_db = Postgresql('Users database',)
                    auth_cache = Redis("auth service cache")

                    scheduler = Container(
                    name="Scheduler Service",
                    technology="Golang + gRPC",
                    description="Scheduler Service responsible for scheduling jobs",)
                    scheduler_db = Mongodb('jobs database')
                    kafka = Kafka("kafka events")

                    email = Container(
                    name="Email Service",
                    technology="Python + gRPC",
                    description="Email Service responsible for executing email jobs",)

                    report = Container(
                    name="Report Service",
                    technology="Python + gRPC",
                    description="Report Service responsible for creating reports",)
                    report_db = Mongodb('reports database')
                
                person >> client
                client >> api_gateway
                api_gateway >> auth
                auth >> auth_db
                auth >> auth_cache
                api_gateway >> scheduler
                scheduler >> scheduler_db
                scheduler >> kafka
                api_gateway >> email
                email >> emaiSystem
                api_gateway >> report
                report >> report_db
                report >> kafka
