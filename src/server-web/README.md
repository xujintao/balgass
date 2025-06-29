# server-web

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

7 Create mirgrations

```
python manage.py makemigrations app1
```

8 Apply mirgrations to migrate database

```
python manage.py migrate
```

9 Run

```
python manage.py runserver 0.0.0.0:8000
```

## Develop

```
./run.sh
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
            "envFile": "${workspaceFolder}/../../config/server-web/.env",
            "django": true,
            "justMyCode": false,
            "console": "integratedTerminal",
        }
    ]
}
```
