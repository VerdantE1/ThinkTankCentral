@echo off
set /p port=Port Number need to be Closed:
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :%port%') do (
    echo Closing the PortNumber %port% of Processor %%a
    taskkill /PID %%a /F
)
echo Port %port% is Closed.
pause