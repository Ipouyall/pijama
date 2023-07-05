from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone

class Disease(models.Model):
    disease_name = models.CharField(max_length=200)
    def __str__(self):
        return self.disease_name

class Hospital(models.Model):
    hospital_name = models.CharField(max_length=200)
    hospital_address = models.CharField(max_length=200)
    def __str__(self):
        return self.hospital_name
    
    