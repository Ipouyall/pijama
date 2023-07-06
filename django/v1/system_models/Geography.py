from django.db import models

class City(models.Model):
    city_name = models.CharField(max_length=30)
    def __str__(self) -> str:
        return self.city_name
