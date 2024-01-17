# server_web

## Prepare

1 Check python3 version

```
python3 --version
# Python 3.10.12
```

2 Create python virtual environment

```
python3 -m venv venv
```

3 Install django

```
pip install django
pip install black
```

4 Freeze requirements.txt

```
pip freeze > requirements.txt
```

5 Create project

```
django-admin startproject project .
```

6 Create app

```
python manage.py startapp app1
```

## Develop

1 Create mirgrations

```
python manage.py makemigrations app1
```

2 Apply mirgrations to migrate database

```
python manage.py migrate
```

3 Run

```
python manage.py runserver 0.0.0.0:8000
```
