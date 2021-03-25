import time
from locust import HttpUser, task, between

class QuickstartUser(HttpUser):
    wait_time = between(1, 2.5)

    @task(3)
    def metrics(self):
        self.client.get("/metrics")

