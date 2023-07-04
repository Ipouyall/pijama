from django.db import models
from django.contrib.auth.models import User
# Create your models here.
from datetime import datetime

class ExtendedUser(models.Model):
    related_user = models.OneToOneField(User,on_delete=models.CASCADE,related_name='extended_user_real_user')
    token = models.CharField(max_length=32,null=True,blank=True)
    
class Disease(models.Model):
    disease_name = models.CharField(max_length=200)
    def __str__(self):
        return self.disease_name

class Hospital(models.Model):
    hospital_name = models.CharField(max_length=200)
    hospital_address = models.CharField(max_length=200)
    def __str__(self):
        return self.hospital_name

class Doctor(models.Model):
    related_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='user_doc')
    practice = models.CharField(max_length=200)
    years_of_experience = models.IntegerField(default=0)
    related_hospitals = models.ManyToManyField(Hospital,verbose_name='hospitals',related_name='hospital_doc')
    def __str__(self):
        return self.related_user.username

class Patient(models.Model):
    related_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='user_patient')
    def __str__(self):
        return self.related_user.username

class Support(models.Model):
    cor_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='support_doc')
    years_of_experience = models.IntegerField(default=0)
    def __str__(self):
        return self.cor_user.username

class City(models.Model):
    city_name = models.CharField(max_length=30)
    def __str__(self) -> str:
        return self.city_name
class Requirement(models.Model):
    description = models.CharField(max_length=200)
    name        = models.CharField(max_length=200)
    def __str__(self):
        return self.name
    
class Package(models.Model):
    package_name = models.CharField(max_length=200)
    description = models.CharField(max_length=200)
    estimate_cost =models.DecimalField(decimal_places=0,max_digits=12)
    requirements = models.ManyToManyField(Requirement,verbose_name='package_reqs')
    city  =models.OneToOneField(City,on_delete=models.DO_NOTHING,null=True,blank=True)
    related_doctor = models.ForeignKey(Doctor,verbose_name='related_doctor',on_delete=models.DO_NOTHING,null=True,blank=True )
    related_hospital = models.ForeignKey(Hospital,verbose_name='related_hospital',on_delete=models.DO_NOTHING,null=True,blank=True )
    disease = models.ForeignKey(Disease,verbose_name='disease',on_delete=models.DO_NOTHING,null=True,blank=True)
    def __str__(self):
        return self.package_name    



class Document(models.Model):
    document_title =models.CharField(max_length=200) 
    content   = models.CharField(max_length=10000,null=True,blank=True)
    related_requirement = models.ForeignKey(Requirement, verbose_name="related_reqs", on_delete=models.DO_NOTHING,
                                            null=True,blank=True)
    def __str__(self):
        return self.document_title

class TreatmentRequestStatus(models.Model):
    status = models.CharField(max_length=200)
    def __str__(self):
        return self.status

class PaymentStatus(models.Model):
    status = models.CharField(max_length=300)
    def __str__(self):
        return self.status
    
class VizaStatus(models.Model):
    status = models.CharField(max_length=300)
    def __str__(self):
        return self.status

class PaymentRequest(models.Model):
    created_date = models.DateTimeField()
    verified_date = models.DateTimeField()
    value =  models.DecimalField(decimal_places=0,max_digits=12)
    related_treatment_request=models.ForeignKey('v1.TreatmentRequest',on_delete=models.CASCADE)
    description = models.CharField(max_length=300)
    status = models.OneToOneField(PaymentStatus,default=1,on_delete=models.DO_NOTHING)

class Viza(models.Model):
    related_user  = models.ForeignKey(User,verbose_name="related_user",on_delete=models.CASCADE)
    expiry_date = models.DateTimeField()
    assigned_date = models.DateTimeField()
    status = models.OneToOneField(VizaStatus,default=1,on_delete=models.DO_NOTHING)
    related_payment_request = models.OneToOneField(PaymentRequest,verbose_name="related_payment_request",on_delete=models.CASCADE,null=True,blank=True)

class TreatmentRequest(models.Model):
    related_package= models.ForeignKey(Package,verbose_name="tr_related_package",on_delete=models.DO_NOTHING,null=True,blank=True)
    related_patient = models.ForeignKey(Patient,verbose_name="tr_related_user",on_delete=models.CASCADE,null=True,blank=True)
    related_documents = models.ManyToManyField(Document,verbose_name="related_documents",null=True,blank=True)
    related_viza  =models.OneToOneField(Viza,null=True,blank=True,on_delete=models.CASCADE)
    submitted_date = models.DateTimeField(auto_created=True)
    last_updated = models.DateTimeField()
    status = models.OneToOneField(TreatmentRequestStatus,default=1,on_delete=models.DO_NOTHING)