from django.db import models
from .Geography import City

class HotelClass(models.Model):
    hotel_class =models.CharField(max_length=300)
    def __str__(self):
        return self.hotel_class
class Hotel(models.Model):
    hotel_name = models.CharField(max_length=300)
    hotel_class = models.ForeignKey(HotelClass,blank=True,on_delete=models.DO_NOTHING,null=True)
    address = models.CharField(max_length=300)
    capacity =models.DecimalField(decimal_places=0,max_digits=12)
    cost =models.DecimalField(decimal_places=0,max_digits=12)
    city = models.ForeignKey(City,blank=True,on_delete=models.DO_NOTHING,null=True)
    def __str__(self):
        return self.hotel_name    
