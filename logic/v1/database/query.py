from ..system_models.ExtendedUser import ExtendedUser, SysAdmin
from ..system_models.Geography import *
from ..system_models.Payment import * 
from ..system_models.Medical import * 
from ..system_models.TreatmentRequest import * 
from ..system_models.Viza import * 
from ..system_models.Requirement import VizaRequirement,Requirement
from ..system_models.Document import VizaDocument
import logging
class QueryBuilder():
    @staticmethod
    def get_user_id_by_payment_id(payment_id):
        return PaymentRequest.objects.filter(pk =payment_id).first().related_treatment_request.related_patient.related_user
    def get_relative_supports(city_name):
        return Support.objects.filter(related_user__city_id=city_name,occupied=False).values_list('id',flat=True)
    def update_support_occupied(support_id):
        support = Support.objects.get(id = support_id)
        support.occupied=True
        support.save()
        return True 
    def get_treatment_request(tr_id,user_id=None):
        if (user_id == None):
            return TreatmentRequest.objects.filter(id=tr_id).first()
        else:
            return TreatmentRequest.objects.filter(id=tr_id,related_patient__id=user_id).first()
            
    def insert_treatment_request(pid,uid,td_ids):
        last_updated = datetime.now()
        ntr = TreatmentRequest(related_package_id =pid,related_patient_id= uid,last_updated=last_updated)
        ntr.save()
        for td_id in td_ids:
            ntr.related_documents.add(td_id)
        ntr.save()
        return ntr.id
    #-----------------------------------------------------------#
    def insert_docs(document_title,document_content,related_requirement_id):
        nd = Document(document_title=document_title,content=document_content,related_requirement_id=related_requirement_id)
        nd.save()
        return nd.id
    def insert_visa_doc(document_title,document_content,related_requirement_id,related_visa_id):
        nd = VizaDocument(related_visa_id=related_visa_id,document_title=document_title,content=document_content,related_requirement_id=related_requirement_id)
        nd.save()
        return nd.id
    #-----------------------------------------------------------#
    def get_hotels(city_id):
        return Hotel.objects.filter(city__id=city_id)
    #-----------------------------------------------------------#
    def get_city(city_name):
        return City.objects.filter(city_name=city_name)
    def get_hotel(hotel_id):
        return Hotel.objects.filter(id=hotel_id).first()
    #-----------------------------------------------------------#
    def get_all_packages():
        return Package.objects.select_related('related_doctor', 'related_hospital', 'disease').all()
    def get_package(id):
        return  Package.objects.filter(pk=id).first()
    
    #-----------------------------------------------------------#
    def get_user_by_token(token):
        return ExtendedUser.objects.filter(token = token).first()
    
    def get_patient_by_token(token):
        return Patient.objects.filter(related_user__token = token).first()
    
    def get_sys_admin_by_token(token):
        return SysAdmin.objects.filter(related_user__token = token).first()
    
    #-----------------------------------------------------------#
    def get_user_by_id(uid):
        return ExtendedUser.objects.filter(related_user__id = uid).first()
    
    #-----------------------------------------------------------#
    def get_sys_admins_chat_id():
        return SysAdmin.objects.all().values_list('related_user__chat_id').all()
    #-----------------------------------------------------------#
    def get_visa_requirements():
        return VizaRequirement.objects.all()
    def get_visa(serial_no,user_token):
        return Viza.objects.filter(serial_no=serial_no,related_user__related_user__token=user_token).first()
    #-----------------------------------------------------------#
    def get_visa_raw(serial_no):
        return Viza.objects.filter(pk=serial_no).first()
    def update_visa_raw_with_new_status(serial_no,new_status):
        return Viza.objects.filter(pk=serial_no).update(status_id=new_status)

    #-----------------------------------------------------------#
    def insert_new_payment_request(value,tr_id):
        p = PaymentRequest(value=value,related_treatment_request_id=tr_id)
        p.save()
        return p
        #-----------------------------------------------------------#
    def get_user_in_tr_id(uid):
        return TreatmentRequest.objects.filter(related_patient__related_user__id=uid).all()
    def get_payment_request(paymentId):
        return PaymentRequest.objects.filter(pk=paymentId).first()
    def update_payment_request_with_status(paymentId,new_status):
        return PaymentRequest.objects.filter(pk=paymentId).update(status_id=new_status)

    def get_user_in_visa(uid):
        viza = Viza.objects.filter(related_user__id=uid).first()
        if (viza == None):
            return viza
        pateint = Patient.objects.filter(related_user__related_user__id=viza.related_user.id).first()
        return pateint