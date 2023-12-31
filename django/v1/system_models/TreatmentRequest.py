from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime
from django.utils import timezone
from .Medical import Disease,Hospital
from .ExtendedUser import Doctor,Patient,Support
from .Geography import City
from .Viza import Viza
from .Hotel import Hotel
from .Requirement import Requirement
from .Document import Document

class PackageClass(models.Model):
    p_class = models.CharField(max_length=200)
    def __str__(self):
        return self.p_class

class Package(models.Model):
    package_name = models.CharField(max_length=200)
    description = models.CharField(max_length=200)
    estimated_cost =models.DecimalField(decimal_places=0,max_digits=12)
    requirements = models.ManyToManyField(Requirement,verbose_name='package_reqs')
    city  =models.ForeignKey(City,on_delete=models.DO_NOTHING,null=True,blank=True)
    related_doctor = models.ForeignKey(Doctor,verbose_name='related_doctor',on_delete=models.DO_NOTHING,null=True,blank=True )
    package_class = models.ForeignKey(PackageClass,verbose_name='package_class',null=True,blank=True,on_delete=models.DO_NOTHING)
    related_hospital = models.ForeignKey(Hospital,verbose_name='related_hospital',on_delete=models.DO_NOTHING,null=True,blank=True )
    disease = models.ForeignKey(Disease,verbose_name='disease',on_delete=models.DO_NOTHING,null=True,blank=True)
    def __str__(self):
        return self.package_name    

TR_STATUS_VERIFYING=1
TR_STATUS_ACTIVE=3
TR_STATUS_DONE=2

class TreatmentRequestStatus(models.Model):
    status = models.CharField(max_length=200)
    def __str__(self):
        return self.status

class TreatmentRequest(models.Model):
    related_package= models.ForeignKey(Package,verbose_name="tr_related_package",on_delete=models.DO_NOTHING,null=True,blank=True)
    related_patient = models.ForeignKey(Patient,verbose_name="tr_related_user",on_delete=models.CASCADE,null=True,blank=True)
    related_documents = models.ManyToManyField(Document,verbose_name="related_documents",null=True,blank=True)
    related_support = models.ForeignKey(Support,verbose_name="related_support",null=True,blank=True,on_delete=models.DO_NOTHING)
    submitted_date = models.DateTimeField(default=timezone.now)
    last_updated = models.DateTimeField()
    status = models.ForeignKey(TreatmentRequestStatus,default=1,on_delete=models.DO_NOTHING)
    reserved_hotel = models.ForeignKey(Hotel,null=True,blank=True,on_delete=models.CASCADE)
    def __str__(self):
        return self.related_package.package_name + ' ' + self.related_patient.related_user.related_user.username
