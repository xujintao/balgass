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
# install
pip install django
pip install black

# or restore:
pip install -r requirements.txt
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

## Debug

### .vscode/launch.json

```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Python: Django",
            "type": "debugpy",
            "request": "launch",
            "program": "${workspaceFolder}/manage.py",
            "args": [
                "runserver",
                "0.0.0.0:8000"
            ],
            "env": {
                "DB_NAME": "django",
                "DB_USER": "root",
                "DB_PASSWORD": "1234",
                "DB_HOST": "localhost",
                "DB_PORT": "5432",
                "TZ": "UTC",
                "EMAIL_HOST": "mail.r2f2.com",
                "EMAIL_PORT": "587",
                "EMAIL_USE_TLS": "1",
                "EMAIL_USER": "noreply@r2f2.com",
                "EMAIL_PASSWORD": "1234",
            },
            "django": true,
            "justMyCode": false,
            "console": "integratedTerminal",
        }
    ]
}
```
