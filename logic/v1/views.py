from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
from django.core import serializers
from v1.models import TreatmentRequest,ExtendedUser,Package,Requirement,Document,Disease,Doctor
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
from django.core.mail import send_mail
TELEGRAM_BOT_API_TOKEN = "6316780921:AAHvZw68iEUCOaTPmRunibA3GyH9--jlbdY"
logger = logging.getLogger(__name__)
class Notifier():
    @staticmethod
    def notify(user_id,message):
        bot =telegram.Bot(token=TELEGRAM_BOT_API_TOKEN)
        # user_chat_id = 
        bot.send_message(user_chat_id=user_id,text=message) 

class QueryBuilder():
    @staticmethod
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
    #-----------------------------------------------------------#
    def get_user_by_token(token):
        return ExtendedUser.objects.filter(token = token).first()

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
    def get_treatment_requests():
        return

class Controller():
    @staticmethod
    def upload_user_docs(request):
        if (request.method == 'POST'):
            request_json = json.loads(request.body)
            pid    = request_json["pid"]
            t_user=QueryBuilder.get_user_by_token(request_json["token"])
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
class PackageHandler():
    @staticmethod
    def get_packages(request):
        packages = list(Package.objects.select_related('related_doctor', 'related_hospital', 'disease').all())
        serialized_packages = json.dumps(list(packages),cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_packages) ,safe=False)

    @staticmethod
    def get_package(request):
        request_json = json.loads(request.body)
        id = request_json["id"]
        package = Package.objects.filter(pk=id)[0]
        serialized_package = json.dumps(package,cls=ExtendedEncoder)
        return JsonResponse(json.loads(serialized_package),safe=False)
    @staticmethod
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

