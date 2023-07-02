from django.contrib import admin
from v1.models import Package,Document,Requirement
# Register your models here.

admin.site.register(Package)
admin.site.register(Document)
admin.site.register(Requirement)
