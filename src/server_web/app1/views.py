from django.shortcuts import render, redirect
from django.contrib.auth import login
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth.decorators import login_required
from django.http import Http404
from . import models
import requests
import json

# Create your views here.


def home(request):
    from django.contrib.auth.models import User

    users = User.objects.all()
    topics = models.Topic.objects.all()
    entries = models.Entry.objects.all()
    context = {"users": users, "topics": topics, "entries": entries}
    return render(request, "home.html", context)


def signup(request):
    if request.method != "POST":
        form = UserCreationForm()
    else:
        form = UserCreationForm(data=request.POST)
        if form.is_valid():
            new_user = form.save()
            login(request, new_user)
            return redirect("home")
    context = {"form": form}
    return render(request, "registration/signup.html", context)


@login_required
def user(request, name):
    context = {}
    if request.user.username == name:
        context["self"] = True
    else:
        try:
            user = models.User.objects.get(username=name)
        except models.User.DoesNotExist:
            return render(request, "404.html")
        else:
            context["self"] = False
            context["user"] = user
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
    return render(request, "game.html")


@login_required
def game_accounts(request):
    context = {}
    # launch api request
    url = "http://r2f2.com:8080/api/accounts"
    try:
        response = requests.get(url, params={"user_id": request.user.id})
    except Exception as e:
        print(e)
        context["get_account_list_message"] = "request server failed"
    else:
        result = response.json()
        if response.status_code == 200:
            context["accounts"] = result
        else:
            context["get_account_list_message"] = result["message"]
    if request.method != "POST":
        form = models.AccountForm()
    else:
        # todo
        form = models.AccountForm(data=request.POST)
        if form.is_valid():
            acc = form.cleaned_data
            headers = {"Content-type": "application/json"}
            param = {
                "name": acc["name"],
                "password": acc["password1"],
                "user_id": request.user.id,
            }
            data = json.dumps(param)
            try:
                response = requests.post(url, data=data, headers=headers)
            except Exception as e:
                print(e)
                context["create_account_message"] = "request server failed"
            else:
                result = response.json()
                if response.status_code != 200:
                    context["create_account_message"] = result["message"]
                else:
                    return redirect(".")
    context["form"] = form
    return render(request, "accounts.html", context)
