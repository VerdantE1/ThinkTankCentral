@echo off
set /p port=Port Number:
netstat -ano | findstr :%port%
if %errorlevel%==0 (
    echo Port %port% Is Using.
) else (
    echo Port %port% Is Not Using.
) 
pause