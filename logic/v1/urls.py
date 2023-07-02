from django.urls import path
from .views import login,get_packages,logout_view,get_package,get_package_requirements

urlpatterns = [
    path("login",login),
    path("logout",logout_view),
    path("packages",get_packages ,name = "get_packages"),
    path("package",get_package,name="get_package"),
    path("package_requirements",get_package_requirements,name="package_requirements")
]
