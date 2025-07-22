from django.shortcuts import render, redirect
from django.contrib.auth import login
from django.contrib.auth.decorators import login_required
from django.contrib.auth.tokens import default_token_generator
from django.conf import settings
from django.http import Http404
from . import models
import requests
import json
import logging

logger = logging.getLogger(__name__)


# Create your views here.
def home(request):
    users = models.User.objects.all()
    topics = models.Topic.objects.all()
    entries = models.Entry.objects.all()
    context = {"users": users, "topics": topics, "entries": entries}
    return render(request, "home.html", context)


def get_client_ip(request):
    x_forwarded_for = request.META.get("HTTP_X_FORWARDED_FOR")
    if x_forwarded_for:
        return x_forwarded_for.split(",")[0].strip()
    return request.META.get("REMOTE_ADDR")


def signup(request):
    context = {"turnstile_sign_up_site_key": settings.TURNSTILE_SIGN_UP_SITE_KEY}
    if request.method != "POST":
        form = models.CustomUserCreationForm()
    else:
        form = models.CustomUserCreationForm(data=request.POST)
        token = request.POST.get("turnstile-response")
        if token:
            verify_url = settings.TURNSTILE_VERIFY_URL
            data = {
                "secret": settings.TURNSTILE_SIGN_UP_SECRET_KEY,
                "response": token,
                "remoteip": get_client_ip(request),
            }
            resp = requests.post(verify_url, data=data)
            result = resp.json()
            if result.get("success"):
                if form.is_valid():
                    new_user = form.save(commit=False)
                    new_user.save()
                    models.Profile.objects.create(user=new_user, email_verified=False)
                    send_verification_email(new_user, request)
                    return render(request, "registration/verification_sent.html")
                    # login(request, new_user)
                    # return redirect("home")
            else:
                form.add_error(None, "CAPTCHA verification failed. Please try again.")
        else:
            form.add_error(None, "Please complete the CAPTCHA.")
    context["form"] = form
    return render(request, "registration/signup.html", context)


def send_verification_email(user, request):
    from django.utils.http import urlsafe_base64_encode
    from django.utils.encoding import force_bytes
    from django.urls import reverse
    from django.core.mail import send_mail

    token = default_token_generator.make_token(user)
    uid = urlsafe_base64_encode(force_bytes(user.pk))
    verification_link = request.build_absolute_uri(
        reverse("verify", kwargs={"uidb64": uid, "token": token})
    )
    subject = "Verify your email address"
    message = f"Click the link to verify your email address: {verification_link}"
    send_mail(subject, message, settings.EMAIL_HOST_USER, [user.email])


@login_required
def resend_verification_email(request):
    if request.method != "POST":
        return redirect("home")
    else:
        user = request.user
        if user.profile.email_verified:
            return redirect("home")
        send_verification_email(user, request)
        return render(request, "registration/verification_sent.html")


def verify(request, uidb64, token):
    from django.utils.http import urlsafe_base64_decode
    from django.utils.encoding import force_str

    try:
        uid = force_str(urlsafe_base64_decode(uidb64))
        user = models.User.objects.get(pk=uid)
    except (TypeError, ValueError, OverflowError, models.User.DoesNotExist):
        user = None
    if user is not None and default_token_generator.check_token(user, token):
        profile = user.profile
        profile.email_verified = True
        profile.save()
        login(request, user)
        return render(request, "registration/verification_success.html")
    else:
        return render(request, "registration/verification_invalid.html")


@login_required
def user(request, name):
    user = request.user
    context = {}
    if user.username != name:
        try:
            other = models.User.objects.get(username=name)
        except models.User.DoesNotExist:
            return render(request, "404.html")
        else:
            context["other"] = other
    return render(request, "user.html", context)


@login_required
def topics(request):
    topics = models.Topic.objects.filter(owner=request.user).order_by("date_added")
    context = {"topics": topics}
    return render(request, "topics.html", context)


@login_required
def topic(request, id):
    topic = models.Topic.objects.get(id=id)
    if topic.owner != request.user:
        raise Http404

    entries = topic.entry_set.order_by("date_added")
    context = {"topic": topic, "entries": entries}
    return render(request, "topic.html", context)


@login_required
def new_topic(request):
    if request.method != "POST":
        form = models.TopicForm()
    else:
        form = models.TopicForm(data=request.POST)
        if form.is_valid():
            new_topic = form.save(commit=False)
            new_topic.owner = request.user
            new_topic.save()
            return redirect("topics")
    context = {"form": form}
    return render(request, "topic_new.html", context)


@login_required
def new_entry(request, id):
    topic = models.Topic.objects.get(id=id)
    if topic.owner != request.user:
        raise Http404

    if request.method != "POST":
        form = models.EntryForm()
    else:
        form = models.EntryForm(data=request.POST)
        if form.is_valid():
            new_entry = form.save(commit=False)
            new_entry.topic = topic
            new_entry.save()
            return redirect("topic", id=id)
    context = {"topic": topic, "form": form}
    return render(request, "entry_new.html", context)


@login_required
def edit_entry(request, id):
    entry = models.Entry.objects.get(id=id)
    topic = entry.topic
    if topic.owner != request.user:
        raise Http404

    if request.method != "POST":
        form = models.EntryForm(instance=entry)
    else:
        form = models.EntryForm(instance=entry, data=request.POST)
        if form.is_valid():
            form.save()
            return redirect("topic", id=topic.id)
    context = {"entry": entry, "topic": topic, "form": form}
    return render(request, "entry_edit.html", context)


def game(request):
    context = {
        "game_websocket_url": settings.GAME_WEBSOCKET_URL,
        "game_map_names": settings.GAME_MAP_NAMES,
    }
    return render(request, "game.html", context)


@login_required
def game_accounts(request):
    if not request.user.profile.email_verified:
        raise Http404
    context = {}
    url = settings.GAME_API_URL
    if request.method != "POST":
        form = models.AccountForm()
    else:
        form = models.AccountForm(data=request.POST)
        if form.is_valid():
            acc = form.cleaned_data
            headers = {"Content-type": "application/json"}
            param = {
                "name": acc["name"],
                "password": acc["password1"],
                "user_email": request.user.email,
            }
            data = json.dumps(param)
            try:
                response = requests.post(url, data=data, headers=headers)
                result = response.json()
                if response.status_code == 200:
                    return redirect(".")
                else:
                    context["create_account_message"] = result["message"]
            except Exception as e:
                logger.error(e)
                context["create_account_message"] = "request server failed"
    try:
        response = requests.get(url, params={"user_email": request.user.email})
        result = response.json()
        if response.status_code == 200:
            context["accounts"] = result
        else:
            context["get_account_list_message"] = result["message"]
    except Exception as e:
        logger.error(e)
        context["get_account_list_message"] = "request server failed"
    context["form"] = form
    return render(request, "accounts.html", context)


def items(request):
    kind = request.GET.get("kind", "sword")
    kind_indexs = settings.GAME_ITEMS_KINDS.get(kind)
    if kind in settings.GAME_ITEMS_STAFF_LIKE_KEYWORDS:
        kind = "staff-like"
    elif kind in settings.GAME_ITEMS_SHIELD_LIKE_KEYWORDS:
        kind = "shield-like"
    elif kind in settings.GAME_ITEMS_WEAPON_KEYWORDS:
        kind = "weapon"
    elif kind in settings.GAME_ITEMS_ARMOR_KEYWORDS:
        kind = "armor"
    elif kind in settings.GAME_ITEMS_WING_KEYWORDS:
        kind = "wing"
    else:
        kind = None
    items = []
    for index in kind_indexs:
        item = settings.GAME_ITEMS.get(index)
        if item is None:
            continue
        items.append(item)
    context = {
        "all_keywords": settings.GAME_ITEMS_KEYWORDS,
        "render_weapon": kind == "weapon",
        "render_staff_like": kind == "staff-like",
        "render_shield_like": kind == "shield-like",
        "render_armor": kind == "armor",
        "render_wing": kind == "wing",
        "items": items,
    }
    return render(request, "items.html", context)
