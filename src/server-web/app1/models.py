from django.db import models
from django.contrib.auth.models import User
from django.contrib.auth.forms import UserCreationForm, AuthenticationForm
from django.contrib.auth import authenticate
from django.core.exceptions import ValidationError
from django import forms


# Create your models here.
class Profile(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    email_verified = models.BooleanField(default=False)


class CustomUserCreationForm(UserCreationForm):
    class Meta(UserCreationForm.Meta):
        fields = (
            "username",
            "email",
        )

    def clean_email(self):
        email = self.cleaned_data.get("email")
        if User.objects.filter(email=email).exists():
            raise ValidationError("Email already exists")
        return email

    def save(self, commit=True):
        user = super().save(commit=False)
        user.email = self.cleaned_data["email"]
        if commit:
            user.save()
            if hasattr(self, "save_m2m"):
                self.save_m2m()
        return user


class CustomAuthenticationForm(AuthenticationForm):
    username = forms.CharField(label="Username or Email")

    def clean(self):
        username = self.cleaned_data.get("username")
        password = self.cleaned_data.get("password")
        if username is not None and password:
            # use username as username
            self.user_cache = authenticate(
                self.request, username=username, password=password
            )
            if self.user_cache is None:
                # use username as email
                try:
                    user = User.objects.get(email=username)
                    self.user_cache = authenticate(
                        self.request, username=user.username, password=password
                    )
                except User.DoesNotExist:
                    pass
        if self.user_cache is None:
            raise self.get_invalid_login_error()
        else:
            self.confirm_login_allowed(self.user_cache)
        return self.cleaned_data


class UserUpdateForm(forms.ModelForm):
    class Meta:
        model = User
        fields = (
            "username",
            "email",
        )

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._old_username = self.instance.username
        self._old_email = self.instance.email

    def clean_username(self):
        username = self.cleaned_data.get("username")
        user_id = self.instance.id
        if User.objects.filter(username=username).exclude(id=user_id).exists():
            raise forms.ValidationError("Username already exists")
        return username

    def clean_email(self):
        email = self.cleaned_data.get("email")
        user_id = self.instance.id
        if User.objects.filter(email=email).exclude(id=user_id).exists():
            raise forms.ValidationError("Email already exists")
        return email

    def save(self, commit=True):
        need_save = False
        need_verify = False
        updated_user = super().save(commit=False)
        if self.cleaned_data["username"] != self._old_username:
            need_save = True
        if self.cleaned_data["email"] != self._old_email:
            need_save = True
            need_verify = True
        if commit and need_save:
            updated_user.save()
        return updated_user, need_verify


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
