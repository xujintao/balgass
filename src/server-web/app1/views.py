from django.shortcuts import render, redirect
from django.contrib.auth import login, update_session_auth_hash
from django.contrib.auth.decorators import login_required
from django.contrib.auth.tokens import default_token_generator
from django.conf import settings
from django.http import Http404
from . import models
from django.contrib.auth.forms import PasswordChangeForm
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
    context["active_url_home"] = True
    return render(request, "home.html", context)


def get_client_ip(request):
    x_forwarded_for = request.META.get("HTTP_X_FORWARDED_FOR")
    if x_forwarded_for:
        return x_forwarded_for.split(",")[0].strip()
    return request.META.get("REMOTE_ADDR")


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
                    return redirect("verification_sent")
                    # login(request, new_user)
                    # return redirect("home")
            else:
                form.add_error(None, "CAPTCHA verification failed. Please try again.")
        else:
            form.add_error(None, "Please complete the CAPTCHA.")
    context["form"] = form
    return render(request, "registration/signup.html", context)


def verification_sent(request):
    user = request.user
    context = {"email": user.email}
    return render(request, "registration/verification_sent.html", context)


from django.contrib.auth.views import LoginView


# login context
class CustomLoginView(LoginView):
    form_class = models.CustomAuthenticationForm

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        context["active_url_login"] = True
        return context


@login_required
def resend_verification_email(request):
    if request.method != "POST":
        return redirect("home")
    else:
        user = request.user
        if user.profile.email_verified:
            return redirect("home")
        send_verification_email(user, request)
        return redirect("verification_sent")


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
def user_settings(request):
    form_update_profile = models.UpdateProfileForm(instance=request.user)
    form_change_password = PasswordChangeForm(user=request.user)
    if request.method == "POST":
        if "update-profile" in request.POST:
            form_update_profile = models.UpdateProfileForm(
                data=request.POST, instance=request.user
            )
            if form_update_profile.is_valid():
                updated_user, need_verify = form_update_profile.save()
                if need_verify:
                    profile = updated_user.profile
                    profile.email_verified = False
                    profile.save()
                    send_verification_email(updated_user, request)
                    return redirect("verification_sent")
                else:
                    return redirect("settings")
        elif "change-password" in request.POST:
            form_change_password = PasswordChangeForm(
                data=request.POST, user=request.user
            )
            if form_change_password.is_valid():
                user = form_change_password.save()
                update_session_auth_hash(request, user)
                return redirect("password_change_done")
    context = {
        "form_update_profile": form_update_profile,
        "form_change_password": form_change_password,
    }
    return render(request, "settings.html", context)


@login_required
def topics(request):
    topics = models.Topic.objects.filter(owner=request.user).order_by("date_added")
    context = {"topics": topics}
    context["active_url_topics"] = True
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
    context["active_url_game"] = True
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
    render_kind = "weapon"
    kind_indexs = settings.GAME_ITEM_KINDS.get(kind)
    if kind in settings.GAME_ITEM_STAFF_LIKE_KEYWORDS:
        render_kind = "staff-like"
    elif kind in settings.GAME_ITEM_SHIELD_KEYWORDS:
        render_kind = "shield"
    elif kind in settings.GAME_ITEM_GLOVES_KEYWORDS:
        render_kind = "gloves"
    elif kind in settings.GAME_ITEM_WING_KEYWORDS:
        render_kind = "wing"
    elif kind in settings.GAME_ITEM_PENDANT_KEYWORDS:
        render_kind = "pendant"
    elif kind in settings.GAME_ITEM_RING_KEYWORDS:
        render_kind = "ring"
    elif kind in settings.GAME_ITEM_SET_KEYWORDS:
        render_kind = "set"
    elif kind in settings.GAME_ITEM_WEAPON_KEYWORDS:
        render_kind = "weapon"
    elif kind in settings.GAME_ITEM_ARMOR_KEYWORDS:
        render_kind = "armor"
    items = []
    if kind == "set":
        sets = settings.GAME_ITEM_SET_INDEXS
        for set in sets:
            set_items = []
            for index in set["item_indexs"]:
                item = settings.GAME_ITEMS.get(index)
                if item is None:
                    continue
                set_items.append(item["name"])
            set["items"] = set_items
            items.append(set)
    else:
        for index in kind_indexs:
            item = settings.GAME_ITEMS.get(index)
            if item is None:
                continue
            items.append(item)
        if render_kind == "wing":
            items = sorted(items, key=lambda item: item["wing_kind"])
        else:
            items = sorted(items, key=lambda item: item["drop_level"])
    context = {
        "all_keywords": settings.GAME_ITEMS_KEYWORDS,
        "kind": kind,
        "render_staff_like": render_kind == "staff-like",
        "render_shield": render_kind == "shield",
        "render_gloves": render_kind == "gloves",
        "render_wing": render_kind == "wing",
        "render_pendant": render_kind == "pendant",
        "render_ring": render_kind == "ring",
        "render_set": render_kind == "set",
        "render_weapon": render_kind == "weapon",
        "render_armor": render_kind == "armor",
        "items": items,
    }
    context["active_url_items"] = True
    return render(request, "items.html", context)
