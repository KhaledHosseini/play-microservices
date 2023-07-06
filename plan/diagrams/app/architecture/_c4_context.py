from diagrams import Diagram, Cluster
from diagrams.c4 import Person,  System, Relationship
class Context:
    def create(self):
        with Diagram("Context", show=False, direction="LR"):
            with Cluster(""):
                person = Person(name='users/admins', description='Any one whos has the application.',external=True)
                app = System(name='microservice app',description='Register users, Schedule email jobs and run them')
                emaiSystem = System(name='emailsystem',description='Runs email-sending jobs')
                person >> Relationship("Register/login, Schedule email jobs to be executed in the future") >> app >> Relationship("Send email jobs to be executed.") >> emaiSystem