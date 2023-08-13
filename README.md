# py-web

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
python manage.py runserver
```

## Deploy
1 Build and push image to docker private registry

```
# git clone
git clone git@github.com:xujintao/py-web.git py-web
cd py-web

# docker build
docker build -t py-web:0.1 .
```

2 Run image

```
docker run \
--restart always \
-d \
--name py-web \
-e TZ=Asia/Shanghai \
-p 8000:8000 \
py-web:0.1
```