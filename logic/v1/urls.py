from django.urls import path
from .views import AccomadationHandler,PackageHandler,AuthenticationHandler,Controller

urlpatterns = [
    path("login",AuthenticationHandler.login),
    path("logout",AuthenticationHandler.logout),
    path("packages",PackageHandler.get_packages ,name = "packages"),
    path("package",PackageHandler.get_package,name="package"),
    path("upload_user_docs",Controller.upload_user_docs,name="upload_user_docs"),
    path("package_requirements",PackageHandler.get_package_requirements,name="package_requirements"),
    path("hotels",Controller.handle_hotel_request,name="hotels"),
]
