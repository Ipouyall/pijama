from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
from django.core import serializers
from .database.query import QueryBuilder
from .system_models.ExtendedUser import ExtendedUser
from .system_models.Geography import *
from .system_models.Payment import * 
from .system_models.Medical import * 
from .system_models.TreatmentRequest import * 
from .system_models.Viza import * 

# from django.contrib.users
from django.contrib.auth import logout                           
import logging
import json
from django.contrib.auth.models import User
from django.contrib.auth import authenticate
from django.contrib.sessions.models import Session
from django.core.serializers.json import DjangoJSONEncoder
from django.db.models import Model   
from v1.encoders import ExtendedEncoder
import random,string
from datetime import datetime
import telegram
import requests
from django.core.mail import send_mail
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
class AccomadationHandler():
    def get_hotels(request):
            request_json = json.loads(request.body)
            package_id = request_json["package_id"]
            package = QueryBuilder.get_package(package_id)
            if (package != None):
                related_package_city_id =package.city.id
                filtered_hotels =list (QueryBuilder.get_hotels(related_package_city_id))
                hotels_json = json.dumps(filtered_hotels,cls=ExtendedEncoder)
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

class Controller():
    def get_user_by_token(request):
        request_json = json.loads(request.body)
        if (request_json.get("token")):
            return QueryBuilder.get_user_by_token(request_json["token"])
        else:
            return None
    @staticmethod
    def upload_user_docs(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            pid    = request_json["pid"]
            t_user=Controller.get_user_by_token(request)
            uid  = 0
            if (t_user != None):
                uid = t_user.related_user.id
                td_ids = DocumentHandler.submit_docs(request_json)
            else:
                return JsonResponse({"status":403,
                                     "message":"Not authenticated"})
                 
            tr_id = TreatmentRequestHandler.create_treatment_request(uid,pid,td_ids)
            return JsonResponse({"tr_id":tr_id,
                                 "status": 200})
        else:
            return JsonResponse({"status":401})
    
    @staticmethod
    def handle_hotel_request(request):
        if (request.method =='GET'):
            return AccomadationHandler.get_hotels(request)
        elif (request.method == 'POST'):
            t_user=Controller.get_user_by_token(request)
            if (t_user !=None):
                    request_json =json.loads(request.body)
                    treatment_request_id = request_json.get("tr_id")
                    treatment_request = TreatmentRequestHandler.get_treatment_request                    (treatment_request_id,t_user.id)
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
        serialized_packages = json.dumps(list(packages),cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_packages) ,safe=False)

    def get_package(request):
        request_json = json.loads(request.body)
        id = request_json["id"]
        package = QueryBuilder.get_package(id)
        serialized_package = json.dumps(package,cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_package),safe=False)
    
    def get_package_requirements(request):
        json_request = json.loads(request.body)
        id = json_request["id"]
        packs = Package.objects.filter(pk=id).first()
        if (packs != None):
            requirements=packs.requirements.all()
        else:
            requirements=[]
        serialized_requirments = json.dumps(list(requirements),cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_requirments) ,safe=False)
    
class AuthenticationHandler():
    @staticmethod
    def logout(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            token  = request_json["token"]
            t_user = QueryBuilder.get_user_by_token(token)
            if (t_user == None):
                return JsonResponse({"message":"Logout failed"
                                    ,"code ": 400})
            else:
                t_user.token = None
                t_user.save()
                return JsonResponse({"message":"Logged out"
                                    ,"code ": 200})
        return HttpResponse(request,"Error Pages/405.html",status = 405)
    @staticmethod
    def login(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            username = request_json["username"]
            password = request_json["password"]
            t_user = ExtendedUser.objects.filter(related_user__username = username).first()
            if (t_user == None):
                return JsonResponse({"message":"User does not exist"
                                    ,"code ": 404})
            else:
                if (not authenticate(username=username,password=password)):
                    return JsonResponse({"message":"Wrong password"
                                        ,"code ": 401})
                else:
                    t_user.token = gen_token()
                    t_user.save()
                    return JsonResponse({"message":"Logged in ",
                                         "token" : t_user.token
                                        ,"code ": 200,})
        return HttpResponse(request,"Error Pages/405.html",status = 405)

