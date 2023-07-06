from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone
from .Payment import PaymentRequest
from .ExtendedUser import Patient
import random,string
Expired_Visa = 2
Active_Visa = 3
Verifying_Visa = 1
class VizaStatus(models.Model):
    status = models.CharField(max_length=300)
    def __str__(self):
        return self.status
# 
class Viza(models.Model):
    related_user  = models.ForeignKey(Patient,verbose_name="related_user",related_name="reverse_user",on_delete=models.CASCADE)
    expiry_date = models.DateTimeField(null=True,blank=True)
    assigned_date = models.DateTimeField(null=True,blank=True)
    serial_no = models.CharField(max_length=64,primary_key=True)
    request_cost = models.DecimalField(max_digits=12,decimal_places=0,default=50000)
    status = models.ForeignKey(VizaStatus,default=Verifying_Visa,on_delete=models.DO_NOTHING,null=True,blank=True)
    related_payment_request = models.OneToOneField(PaymentRequest,verbose_name="related_payment_request",on_delete=models.CASCADE,null=True,blank=True)
    def __str__(self):
        return self.related_user.related_user.related_user.username + "'s Visa"

