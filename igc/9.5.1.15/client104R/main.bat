@echo off
:: Create a temporary drive letter mapped to your UNC root location
:: and effectively CD to that location
pushd %~dp0

:: Do your work
start main.exe connect /u10.0.2.15 /p44405

:: Remove the temporary drive letter and return to your original location
popd
