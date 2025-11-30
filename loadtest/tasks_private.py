from locust import TaskSet, task
import base64
import json
import random
import string

class PrivateTasks(TaskSet):
    def on_start(self):
        u = "admin"
        p = "secret"
        raw = f"{u}:{p}".encode()
        token = base64.b64encode(raw).decode()
        self.auth = {"Authorization": "Basic " + token}

        body = {"msg": "locust-" + "".join(random.choices(string.ascii_letters, k=6))}
        r = self.client.post("/data", headers=self.auth, json=body)
        self.id = r.json()["id"] if r.status_code == 200 else None

    @task
    def get_all(self):
        self.client.get("/data", headers=self.auth)

    @task
    def get_one(self):
        if self.id:
            self.client.get(f"/data/{self.id}", headers=self.auth)

    @task
    def put(self):
        if self.id:
            body = {"msg": "updated"}
            self.client.put(f"/data/{self.id}", headers=self.auth, json=body)

    @task
    def patch(self):
        if self.id:
            body = {"msg": "patched"}
            self.client.patch(f"/data/{self.id}", headers=self.auth, json=body)

    @task
    def delete(self):
        if self.id:
            self.client.delete(f"/data/{self.id}", headers=self.auth)
            body = {"msg": "reset"}
            r = self.client.post("/data", headers=self.auth, json=body)
            self.id = r.json()["id"] if r.status_code == 200 else None
