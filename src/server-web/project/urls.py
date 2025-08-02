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
from django.contrib.auth import views as auth_views
from app1 import views

urlpatterns = [
    path("admin2/", admin.site.urls),
    path("", views.home, name="home"),
    path("login/", views.CustomLoginView.as_view(), name="login"),
    path("logout/", auth_views.LogoutView.as_view(), name="logout"),
    path(
        "password_reset/", auth_views.PasswordResetView.as_view(), name="password_reset"
    ),
    path(
        "password_reset/done/",
        auth_views.PasswordResetDoneView.as_view(),
        name="password_reset_done",
    ),
    path(
        "reset/<uidb64>/<token>/",
        auth_views.PasswordResetConfirmView.as_view(),
        name="password_reset_confirm",
    ),
    path(
        "reset/done/",
        auth_views.PasswordResetCompleteView.as_view(),
        name="password_reset_complete",
    ),
    path("signup/", views.signup, name="signup"),
    path(
        "resend-verification-email/",
        views.resend_verification_email,
        name="resend_verification_email",
    ),
    path("verification-sent", views.verification_sent, name="verification_sent"),
    path("verify/<uidb64>/<token>/", views.verify, name="verify"),
    path("user/<str:name>/", views.user, name="user"),  # profile
    path("settings/", views.user_settings, name="settings"),  # settings
    path(
        "password_change/done/",
        auth_views.PasswordChangeDoneView.as_view(),
        name="password_change_done",
    ),
    path("topics/", views.topics, name="topics"),
    path("topics/<int:id>/", views.topic, name="topic"),
    path("new_topic/", views.new_topic, name="new_topic"),
    path("new_entry/<int:id>/", views.new_entry, name="new_entry"),
    path("edit_entry/<int:id>/", views.edit_entry, name="edit_entry"),
    path("game/", views.game, name="game"),
    path("game/accounts/", views.game_accounts, name="accounts"),
    path("items/", views.items, name="items"),
]
