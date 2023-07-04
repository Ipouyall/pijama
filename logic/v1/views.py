from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
from django.core import serializers
from v1.models import ExtendedUser,Package,Requirement,Document,Disease,Doctor
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
logger = logging.getLogger(__name__)

def gen_token():
    token_length = 32
    return ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(token_length))

    
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
            t_user = ExtendedUser.objects.filter(token = token).first()
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

