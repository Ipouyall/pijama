from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone
from .Payment import PaymentRequest

class VizaStatus(models.Model):
    status = models.CharField(max_length=300)
    def __str__(self):
        return self.status

class Viza(models.Model):
    related_user  = models.ForeignKey(User,verbose_name="related_user",on_delete=models.CASCADE)
    expiry_date = models.DateTimeField()
    assigned_date = models.DateTimeField()
    status = models.ForeignKey(VizaStatus,default=1,on_delete=models.DO_NOTHING)
    related_payment_request = models.OneToOneField(PaymentRequest,verbose_name="related_payment_request",on_delete=models.CASCADE,null=True,blank=True)
    def __str__(self):
        return self.related_user.username + "'s Visa"
