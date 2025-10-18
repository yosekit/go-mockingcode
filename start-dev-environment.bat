@echo off
REM MockingCode Dev Environment Startup Script
echo === Starting MockingCode Development Environment ===
echo.

REM 1. Start Docker Desktop
echo [1/4] Starting Docker Desktop...
if exist "C:\Program Files\Docker\Docker\Docker Desktop.exe" (
    start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    echo Docker Desktop is starting...
    echo Waiting for Docker initialization (30 seconds^)...
    timeout /t 30 /nobreak >nul
) else (
    echo WARNING: Docker Desktop not found
)

REM 2. Start Ubuntu console directly
echo.
echo [2/4] Starting Ubuntu console...
start wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode; exec bash"
timeout /t 2 /nobreak >nul

REM 3. Open Cursor in project directory
echo.
echo [3/4] Launching Cursor in project...
wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode && nohup cursor . > /dev/null 2>&1 &"
timeout /t 3 /nobreak >nul

REM 4. Start docker-compose
echo.
echo [4/4] Starting docker-compose...
wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode && docker-compose -f docker/docker-compose.dev.yml up -d"

echo.
echo === Development Environment Started ===
echo Check the Ubuntu console for confirmation
echo.
echo -----------------------------------------------
echo Type 'exit' and press Enter to stop everything
echo -----------------------------------------------

:WAIT_LOOP
set /p userInput="Enter command: "
if /i "%userInput%"=="exit" goto SHUTDOWN
goto WAIT_LOOP

:SHUTDOWN
echo.
echo === Shutting Down Development Environment ===
echo.

REM 1. Stop Docker containers
echo [1/7] Stopping Docker containers...
wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode && docker-compose -f docker/docker-compose.dev.yml down"
timeout /t 2 /nobreak >nul

REM 2. Close Docker Desktop
echo.
echo [2/7] Closing Docker Desktop...
taskkill /IM "Docker Desktop.exe" /F >nul 2>&1
timeout /t 2 /nobreak >nul

REM 3. Close Ubuntu console (all WSL bash instances)
echo.
echo [3/7] Closing Ubuntu console...
taskkill /IM bash.exe /F >nul 2>&1
timeout /t 1 /nobreak >nul

REM 4. Shutdown WSL
echo.
echo [4/7] Shutting down WSL...
wsl --shutdown
timeout /t 3 /nobreak >nul

REM 5. Check everything is stopped
echo.
echo [5/7] Checking services status...
echo.
echo --- Docker Desktop status:
tasklist /FI "IMAGENAME eq Docker Desktop.exe" 2>nul | find /I /N "Docker Desktop.exe">nul
if "%ERRORLEVEL%"=="0" (
    echo WARNING: Docker Desktop is still running
) else (
    echo OK: Docker Desktop stopped
)

echo.
echo --- WSL status:
wsl -l --running 2>nul | find /I "Ubuntu">nul
if "%ERRORLEVEL%"=="0" (
    echo WARNING: WSL Ubuntu is still running
) else (
    echo OK: WSL stopped
)

REM 6. Close Cursor
echo.
echo [6/7] Closing Cursor...
taskkill /IM cursor.exe /F >nul 2>&1
timeout /t 1 /nobreak >nul

REM 7. Final confirmation
echo.
echo [7/7] Cleanup complete
echo.
echo === Development Environment Stopped ===
echo.
pause


