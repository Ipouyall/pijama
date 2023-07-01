from django.urls import path
from .views import index,login

urlpatterns = [
    path("test", index),
    path("login",login),
]
