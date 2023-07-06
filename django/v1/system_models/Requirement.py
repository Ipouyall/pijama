from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone

class Requirement(models.Model):
    description = models.CharField(max_length=200)
    name        = models.CharField(max_length=200)
    def __str__(self):
        return self.name

class VizaRequirement(Requirement):
    def __str__(self):
        return self.name
