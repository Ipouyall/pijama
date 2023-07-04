from django.urls import path
from .views import PackageHandler,AuthenticationHandler,Controller

urlpatterns = [
    path("login",AuthenticationHandler.login),
    path("logout",AuthenticationHandler.logout),
    path("packages",PackageHandler.get_packages ,name = "get_packages"),
    path("package",PackageHandler.get_package,name="get_package"),
    path("upload_user_docs",Controller.upload_user_docs,name="upload_user_docs"),
    path("package_requirements",PackageHandler.get_package_requirements,name="package_requirements")
]
