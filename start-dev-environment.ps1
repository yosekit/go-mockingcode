# MockingCode Dev Environment Startup Script
# Run from PowerShell with administrator privileges

Write-Host "=== Starting MockingCode Development Environment ===" -ForegroundColor Green

# 1. Start Docker Desktop
Write-Host "`n[1/4] Starting Docker Desktop..." -ForegroundColor Cyan
$dockerPath = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
if (Test-Path $dockerPath) {
    Start-Process $dockerPath
    Write-Host "Docker Desktop is starting..." -ForegroundColor Yellow
    Write-Host "Waiting for Docker initialization (30 seconds)..." -ForegroundColor Yellow
    Start-Sleep -Seconds 30
} else {
    Write-Host "WARNING: Docker Desktop not found at $dockerPath" -ForegroundColor Red
    Write-Host "Please specify the correct path to Docker Desktop" -ForegroundColor Red
}

# 2. Start Ubuntu console directly
Write-Host "`n[2/4] Starting Ubuntu console..." -ForegroundColor Cyan
Start-Process -FilePath "wsl.exe" -ArgumentList "-d", "Ubuntu", "--", "bash", "-c", "cd /usr/projects/mockingcode; exec bash"
Write-Host "Ubuntu console launched" -ForegroundColor Yellow

# 3. Open Cursor in project directory
Write-Host "`n[3/4] Launching Cursor in project..." -ForegroundColor Cyan
Start-Sleep -Seconds 2
try {
    # Launch Cursor via WSL
    wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode && cursor . > /dev/null 2>&1 &"
    Write-Host "Cursor is starting..." -ForegroundColor Yellow
} catch {
    Write-Host "WARNING: Failed to launch Cursor automatically" -ForegroundColor Red
    Write-Host "Open Cursor manually in /usr/projects/mockingcode" -ForegroundColor Red
}

# 4. Start docker-compose
Write-Host "`n[4/4] Starting docker-compose..." -ForegroundColor Cyan
Start-Sleep -Seconds 3
wsl.exe -d Ubuntu -- bash -c "cd /usr/projects/mockingcode && docker-compose -f docker/docker-compose.dev.yml up -d"
Write-Host "Docker containers are starting..." -ForegroundColor Yellow

Write-Host "`n=== Development Environment Started ===" -ForegroundColor Green
Write-Host "Check the Ubuntu console for confirmation" -ForegroundColor Yellow


