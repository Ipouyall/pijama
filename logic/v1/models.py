from django.db import models
# Create your models here.
class Package(models.Model):
    package_name = models.CharField(max_length=200)
    description = models.CharField(max_length=200)
    def __str__(self):
        return self.package_name
    
    
class Requirement(models.Model):
    description = models.CharField(max_length=200)
    name        = models.CharField(max_length=200)
    requirements = models.ManyToManyField(Package,verbose_name='package_reqs')
    def __str__(self):
        return self.name

class Document(models.Model):
    document_title =models.CharField(max_length=200) 
    description = models.CharField(max_length=200)
    corresponding_requirement = models.ForeignKey(Requirement, verbose_name="doc_reqs", on_delete=models.DO_NOTHING)
    corresponding_package     = models.ForeignKey(Package, verbose_name="corres_pack", on_delete=models.DO_NOTHING)
    
    def __str__(self):
        return self.document_title
