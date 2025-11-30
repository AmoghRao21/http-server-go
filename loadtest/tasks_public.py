from locust import TaskSet, task

class PublicTasks(TaskSet):
    @task
    def root(self):
        self.client.get("/")

    @task
    def echo(self):
        self.client.get("/echo?message=hello")

    @task
    def static(self):
        self.client.get("/static/test.txt")
