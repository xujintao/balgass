"""
URL configuration for project project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/4.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path, include
from app1 import views

urlpatterns = [
    path("admin2/", admin.site.urls),
    path("", views.home, name="home"),
    path("", include("django.contrib.auth.urls")),
    path("signup/", views.signup, name="signup"),
    path("user/<str:name>/", views.user, name="user"),
    path("topics/", views.topics, name="topics"),
    path("topics/<int:id>/", views.topic, name="topic"),
    path("new_topic/", views.new_topic, name="new_topic"),
    path("new_entry/<int:id>/", views.new_entry, name="new_entry"),
    path("edit_entry/<int:id>/", views.edit_entry, name="edit_entry"),
    path("game/", views.game, name="game"),
    path("game/accounts/", views.game_accounts, name="accounts"),
]
