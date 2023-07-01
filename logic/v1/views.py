from django.shortcuts import render
from django.http import HttpResponse,JsonResponse
import logging
import json
logger = logging.getLogger(__name__)
                           
def index(request):
    response = {'Congrats':'GHasemi behtarin ostad donyast'}
    return JsonResponse(response)

def login(request):
    request_json = json.loads(request.body)
    logger.warning(request_json) 
    # response = {"Token" : request.body}
    return JsonResponse(request_json)