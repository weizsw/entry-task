import random

from locust import HttpUser, TaskSet, task


class Testlocust(TaskSet):
    username = "test"
    password = "test"

    def on_start(self):
        num = random.randint(1, 1000)
        self.username += str(num)

    @task
    def login(self):
        self.client.post(
            "/login", {"username": self.username, "password": self.password}
        )


class WebsiteUser(HttpUser):
    tasks = [Testlocust]
    min_wait = 100
    max_wait = 1000
