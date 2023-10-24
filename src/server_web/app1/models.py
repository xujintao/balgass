from django.db import models
from django.contrib.auth.models import User
from django import forms

# Create your models here.


class Topic(models.Model):
    text = models.CharField(max_length=200)
    date_added = models.DateTimeField(auto_now_add=True)
    owner = models.ForeignKey(User, on_delete=models.CASCADE)

    def __str__(self) -> str:
        return self.text


class TopicForm(forms.ModelForm):
    class Meta:
        model = Topic
        fields = ["text"]
        labels = {"text": ""}


class Entry(models.Model):
    """specific knowledge of the topic"""

    topic = models.ForeignKey(Topic, on_delete=models.CASCADE)
    text = models.TextField()
    date_added = models.DateTimeField(auto_now_add=True)

    class Meta:
        verbose_name_plural = "entries"

    def __str__(self) -> str:
        return f"{self.text[:50]}..."


class EntryForm(forms.ModelForm):
    class Meta:
        model = Entry
        fields = ["text"]
        labels = {"text": ""}
        widgets = {"text": forms.Textarea(attrs={"cols": 80})}


class AccountForm(forms.Form):
    # name = forms.CharField(label="Account name", max_length=10)
    name = forms.CharField(
        label="account",
        max_length=10,
    )
    password1 = forms.CharField(
        label="password",
        widget=forms.PasswordInput(),
        max_length=10,
    )
    password2 = forms.CharField(
        label="password confirmation",
        widget=forms.PasswordInput(),
        max_length=10,
    )
