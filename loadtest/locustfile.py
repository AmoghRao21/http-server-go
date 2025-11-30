from locust import HttpUser, between
from tasks_public import PublicTasks
from tasks_private import PrivateTasks

class FullUser(HttpUser):
    wait_time = between(0.1, 0.3)

    def on_start(self):
        pass

    tasks = {
        PublicTasks: 50,
        PrivateTasks: 50,
    }
