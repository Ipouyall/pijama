# TODO check city id stuff  DONE
# TODO fix status codes DONE
# TODO Payment Request double endpoint that calls assigns support
# TODO Notify different people
# TODO Move queries of package and hotel to Query Builder

from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
from django.core import serializers
from .database.query import QueryBuilder
from .system_models.Geography import *
from .system_models.Payment import * 
from .system_models.Medical import * 
from .system_models.TreatmentRequest import * 
from .system_models.Viza import * 
from .system_models.ExtendedUser import * 

import random
# from django.contrib.users
from django.contrib.auth import logout                           
import logging
import json
from django.contrib.auth.models import User
from django.contrib.auth import authenticate
from django.contrib.sessions.models import Session
from django.core.serializers.json import DjangoJSONEncoder
from django.db.models import Model   
from v1.encoders import ExtendedEncoder, ModelToDict
import random,string
from datetime import datetime
import requests
# from django.core.mail import send_mail
TELEGRAM_BOT_API_TOKEN = "6316780921:AAHvZw68iEUCOaTPmRunibA3GyH9--jlbdY"
logger = logging.getLogger(__name__)
class Notifier():
    @staticmethod
    def notify(chat_id,message):
       bot_token = '6316780921:AAHvZw68iEUCOaTPmRunibA3GyH9--jlbdY'
       url = f'https://api.telegram.org/bot{bot_token}/sendMessage'
       params = {'chat_id': chat_id, 'text': message}
       return requests.post(url, params=params)


def gen_token(token_length=32):
    return ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(token_length))

class DocumentHandler():
    @staticmethod
    def submit_docs(request_json):
            documents = request_json["documents"] # Map of document contents and requirement id and document title
            td_ids = []
            for document in documents:
                title = document["title"]
                related_req_id = document["related_req_id"]
                content = document ["content"]
                id = QueryBuilder.insert_docs(title,content,related_req_id)
                td_ids.append(id)
            return td_ids

class TreatmentRequestHandler():
    @staticmethod
    def create_treatment_request(pid,uid,td_ids):
        return QueryBuilder.insert_treatment_request(pid,uid,td_ids)
    def get_treatment_request(tr_id,user_id=None):
        return QueryBuilder.get_treatment_request(tr_id,user_id)
    def update_treatment_request_with_hotel_id(tr_id,hotel_id):
        treatment_request = TreatmentRequestHandler.get_treatment_request(tr_id)
        treatment_request.reserved_hotel_id=hotel_id
        treatment_request.save()
        return True
    def update_treatment_request_with_support_id(tr_id,support_id):
        return QueryBuilder.update_treatment_request_with_support_id(tr_id,support_id)

class AccomadationHandler():
    def get_hotels(request):
            request_json = json.loads(request.body)
            package_id = request_json["package_id"]
            package = QueryBuilder.get_package(package_id)
            if (package != None):
                related_package_city_id =package.city.id
                filtered_hotels =list (QueryBuilder.get_hotels(related_package_city_id))

                hotels_dict = ModelToDict(list(filtered_hotels))
                for hotel_dict in hotels_dict:
                    hotel_dict['city'] = hotel_dict['city']['city_name']
                    hotel_dict['hotel_class'] = hotel_dict['hotel_class']['hotel_class']

                hotels_json = json.dumps(hotels_dict,cls=ExtendedEncoder)
                return JsonResponse(json.loads(hotels_json) ,safe=False)
            else:
                response = JsonResponse({"message":"Package does not exist"},safe=False)
                response.status_code = 404 
                return response
            
    def reserve_hotel(request):
            request_json = json.loads(request.body)
            hotel_id = request_json.get("hotel_id")
            hotel = QueryBuilder.get_hotel(hotel_id)
            if (hotel != None):
                hotel.capacity=hotel.capacity -1 
                if (hotel.capacity < 0):
                    return False
                else:
                    hotel.save()
                    return True
            else:
                return False
class VisaHandler():
    def get_visa_requirements():
        serialized_requirments = json.dumps(list(QueryBuilder.get_visa_requirements()),cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_requirments) ,safe=False)
    
    def create_empty_visa(user_id):
        serial_no = gen_token(64)
        new_visa =Viza(serial_no=serial_no,related_user_id=user_id)
        new_visa.save()
        return new_visa
    
    def add_visa_docs(request_json,serial_no):
        documents = request_json["documents"] # Map of document contents and requirement id and document title
        for document in documents:
                        title = document["title"]
                        related_req_id = document["related_req_id"]
                        content = document ["content"]
                        id = QueryBuilder.insert_visa_doc(title,content,related_req_id,serial_no)
        return True
    def check_viza(serial_no,user_token):
        viza= QueryBuilder.get_visa(serial_no,user_token)
        return viza            
    
class PaymentHandler():
    @staticmethod
    def create_invoice(value,tr_id):
        return QueryBuilder.insert_new_payment_request(value,tr_id)

class SupportHandler():
    @staticmethod
    def assign_support(city_name):
        logger.warn(city_name)
        potential_support_ids = QueryBuilder.get_relative_supports(city_name)
        if (potential_support_ids.__len__() != 0):
            random_support_id = random.choices(potential_support_ids)
            return random_support_id[0]
        else:
            return False
    
    def update_support_occupied(support_id):
        support_occupied_id = QueryBuilder.update_support_occupied(support_id)
        return support_occupied_id
class Controller():
    def get_user_id_by_payment_id(payment_id):
        return QueryBuilder.get_user_id_by_payment_id()
    def change_visa_status(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            
            tr_user = Controller.get_user_by_token(request)
            logging.warn(tr_user)
            if (tr_user != None and tr_user.related_user.username == "admin"):
                serial_no = request_json.get("serial_no") 
                visa = QueryBuilder.get_visa_raw(serial_no)
                user = visa.related_user.related_user
                if (user == None):
                    response = JsonResponse({"message":"User does not exist"},safe=False)
                    response.status_code = 404 
                    return response

                if (visa != None):
                    logging.warn(visa.status.id)
                    new_status = request_json.get("visa_status")
                    logging.warn(request_json.get("visa_status"))
                    logging.warn(visa.status.id)
                    QueryBuilder.update_visa_raw_with_new_status(serial_no,new_status)
                    Notifier.notify(user.chat_id, "Your viza accepted.")
                    return JsonResponse({"message":"Visa status updated succefully"},safe=False)
                else:
                    response = JsonResponse({"message":"Visa does not exist"},safe=False)
                    response.status_code = 404 
                    return response
            else:
                response = JsonResponse({"message":"You don't have access."},safe=False)
                response.status_code = 403 
                return response
    @staticmethod
    def change_payment_status(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            
            tr_user = Controller.get_sys_admin_by_token(request)
            if (tr_user != None):
                user = QueryBuilder.get_user_id_by_payment_id(request_json.get("payment_id"))
                extended_user = user.related_user
                if (user == None):
                    response = JsonResponse({"message":"User does not exist"},safe=False)
                    response.status_code = 404 
                    return response

                payment_id = request_json.get("payment_id") 
                payment_status = request_json.get("payment_status")
                payment_request = QueryBuilder.get_payment_request(payment_id)
                treatment_request= payment_request.related_treatment_request
                if (payment_request != None):
                    logging.warn(payment_request.status.id)
                    logging.warn(request_json.get("payment_id"))
                    logging.warn(payment_request.status.id)
                    QueryBuilder.update_payment_request_with_status(payment_id,payment_status)
                    
                    Notifier.notify(extended_user.chat_id, format("Your payment has been accepted by admin"))
                    support_id =Controller.assign_support(treatment_request.id,user.id)
                    if (support_id):
                        TreatmentRequestHandler.update_treatment_request_with_support_id(treatment_request.id,support_id)
                        Notifier.notify(extended_user.chat_id, format("Support has been assigned to your treatment request"))
                    else :
                        Notifier.notify(extended_user.chat_id, format("Failed to assign support to your treatment request please check the panel"))
                        
                    return JsonResponse({"message":"PaymentRequest status updated succefully"},safe=False)
                else:
                    response = JsonResponse({"message":"PaymentRequest does not exist"},safe=False)
                    response.status_code = 404
                    return response
            else:
                response = JsonResponse({"message":"You don't have access."},safe=False)
                response.status_code = 403 
                return response

    @staticmethod
    # Change Treatment Request Status
    def assign_support(tr_id,user_id):
            tr_user = QueryBuilder.get_user_by_id(user_id)
            if (tr_user != None):
                treatment_request = TreatmentRequestHandler.get_treatment_request(tr_id,tr_user.id)
                if (treatment_request != None):
                    assigned_support_id = SupportHandler.assign_support(treatment_request.related_package.city.id)
                    if (assigned_support_id):
                        SupportHandler.update_support_occupied(assigned_support_id)
                        return assigned_support_id
                    else:
                        return False
                else:
                    return False
            else:   
                return False
    def handle_payment_bill_request(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            tr_user = Controller.get_user_by_token(request)
            if (tr_user != None):
                treatment_request = QueryBuilder.get_user_in_tr_id(tr_user.related_user.id).first()
                logging.warn(treatment_request)
                if (treatment_request != None):
                    serial_no = request_json.get("serial_no")
                    visa = VisaHandler.check_viza(serial_no,tr_user.token)
                    if (visa != None):
                        if (visa.status.id == Active_Visa):
                            total_cost  = treatment_request.related_package.estimated_cost + treatment_request.reserved_hotel.cost 
                            package_payment = PaymentHandler.create_invoice(total_cost,treatment_request.id)
                            return JsonResponse({"status":200,
                                                 "message":"Package payment request created",
                                                 "payment_request_id": package_payment.id,
                                                 "total_cost" : int(total_cost) })
                        else:
                            response =JsonResponse({"message":"Visa status is " + visa.status.status + " Wait for viza confirmation or try again for viza"}) 
                            response.status_code = 403
                            return response    
                    else:
                        response =JsonResponse({"message":"Visa did not exist for the provided user"}) 
                        response.status_code = 404
                        return response   
                else:
                 response =JsonResponse({"message":"Not Authorized to access Treatment Request or Treatment Request not found"}) 
                 response.status_code = 403
                 return response 
            else:
                response = JsonResponse({"message":"You don't have access."},safe=False)
                response.status_code = 403 
                return response
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response  
    
    def get_user_by_token(request):
        request_json = json.loads(request.body)
        if (request_json.get("token")):
            return QueryBuilder.get_user_by_token(request_json["token"])
        else:
            return None
    def get_patient_by_token(request):
        request_json = json.loads(request.body)
        if (request_json.get("token")):
            return QueryBuilder.get_patient_by_token(request_json["token"])
        else:
            return None
   
    def get_sys_admin_by_token(request):
        request_json = json.loads(request.body)
        if (request_json.get("token")):
            return QueryBuilder.get_sys_admin_by_token(request_json["token"])
        else:
            return None
   
    def handle_visa_request(request):
        if (request.method =='GET'):
            return VisaHandler.get_visa_requirements()
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            tr_user = Controller.get_patient_by_token(request)
            if (tr_user != None):
                # Check tr id exists 
                tr_id = request_json.get("tr_id")
                treatment_request = TreatmentRequestHandler.get_treatment_request(tr_id,tr_user.id)
                if (treatment_request != None):
                    logger.warn(tr_user.id)
                    visa = VisaHandler.create_empty_visa(tr_user.id)
                    VisaHandler.add_visa_docs(request_json,visa.serial_no)
                    payment_request = PaymentHandler.create_invoice(visa.request_cost,tr_id)
                    visa.related_payment_request_id=payment_request.id
                    visa.save()
                    return JsonResponse({"status":200,
                                         "message":"Visa request pending succesfully,Payment Request created and needs to be paid",
                                         "serial_no":visa.serial_no,
                                         "payment_request_id" :payment_request.id,
                                         "total_cost": payment_request.value })
                else:
                 response =JsonResponse({"message":"Treatment Request doesn't exist"}) 
                 response.status_code = 403
                 return response
            else:
                 response =JsonResponse({"message":"Not Authorized"}) 
                 response.status_code = 403
                 return response  
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response
    def check1(request):
        request_json = json.loads(request.body)



    def get_visa_status(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            serial_no = request_json.get("serial_no")
            tr_user = Controller.get_user_by_token(request)
            if (tr_user != None):
                message = ''
                # Todo potential fix
                http_code = 200
                visa=VisaHandler.check_viza(serial_no,tr_user.token)
                if (visa == None):
                    message = "There was no visa try uploading and applying with the corresponding documents"
                elif (visa.status_id == Expired_Visa):
                    message =  "Your visa has expired try again by applying for another visa"
                elif (visa.status_id == Verifying_Visa):
                    message =  "Your Visa is being verified please wait"
                elif (visa.status_id == Active_Visa):
                    message = "You already have an active visa and don't need to apply for a visa"
                
                response = JsonResponse({"message":message,
                                         "serial_no":serial_no})
                response.status_code==http_code
                return response
            else:
             response = JsonResponse({"status":403,
                                     "message":"Not Authorized"}) 
             response.status_code = 403
             return response  
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response
    def upload_user_docs(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            pid    = request_json.get("pid")
            package = QueryBuilder.get_package(pid)
            if (package == None):
                response = JsonResponse({"status":403,
                                     "message":"PAckage not found"})
                response.status_code = 403
                return response
            t_user = Controller.get_patient_by_token(request)
            uid = 0
            if (t_user != None):
                trs = list(QueryBuilder.get_user_in_tr_id(t_user.related_user.pk))
                logger.warn(trs)
                for tr in trs:
                    logger.warn(tr.status.id)
                    if(tr.status.id == 1):
                        # TODO check status code
                        response = JsonResponse({"status":403,
                                     "message":"You have already have an active treatment request: " + str(tr.id)})
                        response.status_code = 403
                        return response
                uid = t_user.id
                td_ids = DocumentHandler.submit_docs(request_json)
            else:
                response = JsonResponse({"status":403,
                                     "message":"Not Authorized"})
                response.status_code = 403
                return response
                 
            logger.warn([uid,pid,td_ids])
            tr_id = TreatmentRequestHandler.create_treatment_request(pid, uid, td_ids)
            return JsonResponse({"tr_id": tr_id,
                                 "status": 200})
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 401
            return response
    
    def handle_hotel_request(request):
        if (request.method =='GET'):
            return AccomadationHandler.get_hotels(request)
        elif (request.method == 'POST'):
            t_patient=Controller.get_patient_by_token(request)
            if (t_patient !=None):
                    request_json =json.loads(request.body)
                    treatment_request_id = request_json.get("tr_id")
                    treatment_request = TreatmentRequestHandler.get_treatment_request(treatment_request_id,t_patient.id)
                    if (treatment_request == None):
                        response =JsonResponse({"message":"Treatment request was not found"})
                        response.status_code = 403
                        return response
                    else:
                        reserve = AccomadationHandler.reserve_hotel(request)
                        if (reserve):
                            hotel_id=request_json.get("hotel_id")
                            TreatmentRequestHandler.update_treatment_request_with_hotel_id(treatment_request_id,hotel_id)
                            response =JsonResponse({"message":"Reserved and added hotel to treatment request "})
                            return response
                        else:
                            response =JsonResponse({"message":"Hotel was not found or not available for reservation"})
                            response.status_code = 403
                            return response
            else:
                response = JsonResponse({"message":"Unauthorized"})
                response.status_code=401
                return response
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response
        
        
class PackageHandler():
    @staticmethod
    def get_packages(request):
        packages = list(QueryBuilder.get_all_packages())
        packages_dict = ModelToDict(list(packages))
        for package_dict in packages_dict:
            package_dict['city'] = package_dict['city']['city_name']
            package_dict['related_doctor'] = package_dict['related_doctor']['related_user']['related_user']['username']
            package_dict['doctor'] = package_dict.pop('related_doctor')
            package_dict['related_hospital'] = package_dict['related_hospital']['hospital_name']
            package_dict['hospital'] = package_dict.pop('related_hospital')
            package_dict['disease'] = package_dict['disease']['disease_name']
            package_dict['package_class'] = package_dict['package_class']['p_class']
            package_dict['catagory'] = package_dict.pop('disease')
            package_dict['estimated_cost'] = int(package_dict['estimated_cost'] )
            del package_dict['requirements']
        serialized_packages = json.dumps(list(packages_dict),cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_packages) ,safe=False)

    def get_package(request):
        request_json = json.loads(request.body)
        id = request_json["id"]
        package = QueryBuilder.get_package(id)
        package_dict = ModelToDict(package)
        package_dict['city'] = package_dict['city']['city_name']
        package_dict['related_doctor'] = package_dict['related_doctor']['related_user']['related_user']['username']
        package_dict['doctor'] = package_dict.pop('related_doctor')
        package_dict['related_hospital'] = package_dict['related_hospital']['hospital_name']
        package_dict['hospital'] = package_dict.pop('related_hospital')
        package_dict['disease'] = package_dict['disease']['disease_name']
        package_dict['package_class'] = package_dict['package_class']['p_class']
        package_dict['catagory'] = package_dict.pop('disease')
        package_dict['estimated_cost'] = int(package_dict['estimated_cost'] )
        del package_dict['requirements']
        serialized_package = json.dumps(package_dict,cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_package),safe=False)
    
    def get_package_requirements(request):
        json_request = json.loads(request.body)
        id = json_request["id"]
        packs = Package.objects.filter(pk=id).first()
        if (packs != None):
            requirements=packs.requirements.all()
            serialized_requirments = json.dumps(list(requirements),cls=ExtendedEncoder)
            return JsonResponse(json.loads(serialized_requirments) ,safe=False)
        else:
            response = JsonResponse({"message":"Package does not exist"
                                    , "code":404})
            response.status_code=404
            return response
class AuthenticationHandler():
    @staticmethod
    def logout(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            token  = request_json["token"]
            t_user = QueryBuilder.get_user_by_token(token)
            if (t_user == None):
                response = JsonResponse({"message":"Logout failed"
                                    ,"code ": 400})
                response.status_code = 400
                return response
            else:
                t_user.token = None
                t_user.save()
                return JsonResponse({"message":"Logged out"
                                    ,"code ": 200})
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response

    def login(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            username = request_json["username"]
            password = request_json["password"]
            t_user = ExtendedUser.objects.filter(related_user__username = username).first()
            if (t_user == None):
                response = JsonResponse({"message":"User does not exist"
                                    ,"code ": 404})
                response.status_code = 404
                return response
            else:
                if (not authenticate(username=username,password=password)):
                    response = JsonResponse({"message":"Wrong password"
                                        ,"code ": 401})
                    response.status_code = 401
                    return response
                else:
                    t_user.token = gen_token()
                    t_user.save()
                    return JsonResponse({"message":"Logged in ",
                                         "token" : t_user.token
                                        ,"code ": 200,})
        else:
            response = JsonResponse({"message":"Method not allowed"},safe=False)
            response.status_code = 405 
            return response

