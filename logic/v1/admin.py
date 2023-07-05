from django.contrib import admin
from .system_models.ExtendedUser import *
from .system_models.Geography import *
from .system_models.Payment import * 
from .system_models.Hotel import *
from .system_models.Medical import * 
from .system_models.TreatmentRequest import * 
from .system_models.Viza import * 
from .system_models.Document import *
from .system_models.Requirement import *
# Register your models here.

admin.site.register(Package)
admin.site.register(Document)
admin.site.register(Requirement)
admin.site.register(Disease)
admin.site.register(Doctor)
admin.site.register(City)
admin.site.register(PaymentStatus)
admin.site.register(VizaStatus)
admin.site.register(TreatmentRequestStatus)
admin.site.register(TreatmentRequest)
admin.site.register(Viza)
admin.site.register(Hospital)
admin.site.register(PaymentRequest)
admin.site.register(ExtendedUser)
admin.site.register(Patient)
admin.site.register(Hotel)
admin.site.register(HotelClass)
admin.site.register(PackageClass)
admin.site.register(SysAdmin)
admin.site.register(VizaRequirement)
admin.site.register(VizaDocument)