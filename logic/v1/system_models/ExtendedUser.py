from django.db import models
from django.contrib.auth.models import User
from .Medical import Hospital

class ExtendedUser(models.Model):
    related_user = models.OneToOneField(User,on_delete=models.CASCADE,related_name='extended_user_real_user')
    chat_id =models.IntegerField()
    token = models.CharField(max_length=32,null=True,blank=True)
    def __str__(self) -> str:
        return self.related_user.username
        
class Doctor(models.Model):
    related_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='user_doc')
    practice = models.CharField(max_length=200)
    years_of_experience = models.IntegerField(default=0)
    related_hospitals = models.ManyToManyField(Hospital,verbose_name='hospitals',related_name='hospital_doc')
    def __str__(self):
        return self.related_user.related_user.username

class Patient(models.Model):
    related_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='user_patient')
    def __str__(self):
        return self.related_user.related_user.username

class Support(models.Model):
    cor_user = models.OneToOneField(ExtendedUser,verbose_name='related_user',on_delete=models.CASCADE,related_name='support_doc')
    years_of_experience = models.IntegerField(default=0)
    def __str__(self):
        return self.related_user.related_user.username
