from django.urls import path
from .views import login,logout_view,PackageHandler

urlpatterns = [
    path("login",login),
    path("logout",logout_view),
    path("packages",PackageHandler.get_packages ,name = "get_packages"),
    path("package",PackageHandler.get_package,name="get_package"),
    path("package_requirements",PackageHandler.get_package_requirements,name="package_requirements")
]
