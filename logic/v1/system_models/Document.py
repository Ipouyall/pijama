
from .Requirement import Requirement,VizaRequirement
from .Viza import Viza
from django.db import models


class Document(models.Model):
    document_title =models.CharField(max_length=200) 
    content   = models.CharField(max_length=10000,null=True,blank=True)
    related_requirement = models.ForeignKey(Requirement, verbose_name="related_reqs", on_delete=models.DO_NOTHING,
                                            null=True,blank=True)
    def __str__(self):
        return self.document_title


class VizaDocument(Document):
    related_visa = models.ForeignKey(Viza,verbose_name="related_visa",on_delete=models.DO_NOTHING,null=True,blank=True) 
    def __str__(self):
        return self.document_title
