from django.shortcuts import render, redirect
from django.contrib.auth import login
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth.decorators import login_required
from django.http import Http404
from . import models

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
