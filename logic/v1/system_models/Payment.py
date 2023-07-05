from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone

class PaymentStatus(models.Model):
    status = models.CharField(max_length=300)
    def __str__(self):
        return self.status

class PaymentRequest(models.Model):
    created_date = models.DateTimeField(default=timezone.now)
    verified_date = models.DateTimeField()
    value =  models.DecimalField(decimal_places=0,max_digits=12)
    related_treatment_request=models.ForeignKey('v1.TreatmentRequest',on_delete=models.CASCADE)
    description = models.CharField(max_length=300)
    status = models.ForeignKey(PaymentStatus,default=1,on_delete=models.DO_NOTHING)
    def __str__(self):
        return self.id
