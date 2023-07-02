from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
from django.core import serializers
from v1.models import Package,Requirement,Document
# from django.contrib.users
import logging
import json
from django.core.serializers.json import DjangoJSONEncoder
from django.contrib.auth.models import User
from django.contrib.auth import authenticate
from django.contrib.sessions.models import Session
from django.forms import model_to_dict
from django.core.serializers.json import DjangoJSONEncoder
from django.db.models import Model

class ExtendedEncoder(DjangoJSONEncoder):
    def default(self, o):
        if isinstance(o, Model):
            return model_to_dict(o)
        return super().default(o)
    
logger = logging.getLogger(__name__)
from django.contrib.auth import logout                           

def logout_view(request):
    if (request.method == 'POST'):
        request_json = json.loads(request.body)
        username = request_json["username"]
        password = request_json["password"]
        t_user = User.objects.filter(username = username)[0]
        real_password = t_user.password
        if (real_password == None):
            return JsonResponse({"message":"User does not exist"
                                 ,"code ": 404})
        elif (not t_user.is_authenticated):
            return JsonResponse({"message":"Not logged in"
                                 ,"code ": 401})
        else:
            user = User.objects.get(username=username)
            [s.delete() for s in Session.objects.all() if s.get_decoded().get('_auth_user_id') == user.id]
            user.save()
            return JsonResponse({"message":"Logged out"
                                 ,"code ": 200})
    return HttpResponse(request,"Error Pages/405.html",status = 405)

def login(request):
    if (request.method == 'POST'):
        request_json = json.loads(request.body)
        username = request_json["username"]
        password = request_json["password"]
        t_user = User.objects.filter(username = username)[0]
        real_password = t_user.password
        if (real_password == None):
            return JsonResponse({"message":"User does not exist"
                                 ,"code ": 404})
        elif (not authenticate(username=username,password=password)):
            return JsonResponse({"message":"Wrong password"
                                 ,"code ": 401})
        else:
            t_user.save()
            return JsonResponse({"message":"Logged in "
                                 ,"code ": 401,})
    return HttpResponse(request,"Error Pages/405.html",status = 405)

def get_packages(request):
    # if (request.user.is_authenticated):
    packages = Package.objects.all()
    serialized_packages = serializers.serialize("json",packages)
    return JsonResponse(json.loads(serialized_packages) ,safe=False)

def get_package(request):
    request_json = json.loads(request.body)
    id = request_json["id"]
    package = Package.objects.filter(pk=id)[0]
    serialized_package = json.dumps(package,cls=ExtendedEncoder)
    return JsonResponse(json.loads(serialized_package),safe=False)

def get_package_requirements(request):
    request_json = json.loads(request.body)
    id = request_json["id"]
    reqs = Requirement.objects.filter(requirements=id)
    serialized_packages = serializers.serialize("json",reqs)
    return JsonResponse(json.loads(serialized_packages) ,safe=False)
