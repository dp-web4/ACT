@echo off
REM Test script for the event system (Windows version)
REM This script demonstrates how to test the event emission system

echo === Event System Test ===
echo This script will:
echo 1. Create a test configuration with events enabled
echo 2. Build and start the API bridge
echo 3. Test component registration with event emission
echo 4. Show expected events
echo.

REM Check if required tools are available
echo [INFO] Checking requirements...

where curl >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] curl is required but not installed
    echo Please install curl from https://curl.se/windows/
    pause
    exit /b 1
)

echo [INFO] Requirements check passed

REM Create test configuration with events enabled
echo [INFO] Creating test configuration with events enabled...

(
echo blockchain:
echo   rest_endpoint: "http://0.0.0.0:1317"
echo   grpc_endpoint: "localhost:9090"
echo   chain_id: "racecarweb"
echo   timeout: 30
echo.
echo server:
echo   port: 8081
echo   host: "0.0.0.0"
echo   read_timeout: 30
echo   write_timeout: 30
echo.
echo logging:
echo   level: "info"
echo   format: "json"
echo.
echo events:
echo   enabled: true
echo   max_retries: 3
echo   retry_delay: 2
echo   queue_size: 1000
echo   endpoints:
echo     component_registered:
echo       - "http://localhost:3000/webhooks/component-registered"
echo     component_verified:
echo       - "http://localhost:3000/webhooks/component-verified"
echo     pairing_initiated:
echo       - "http://localhost:3000/webhooks/pairing-initiated"
echo     pairing_completed:
echo       - "http://localhost:3000/webhooks/pairing-completed"
echo     lct_created:
echo       - "http://localhost:3000/webhooks/lct-created"
echo     trust_tensor_created:
echo       - "http://localhost:3000/webhooks/trust-tensor-created"
echo     energy_transfer:
echo       - "http://localhost:3000/webhooks/energy-transfer"
) > test_config.yaml

echo [INFO] Test configuration created: test_config.yaml

REM Build the API bridge
echo [INFO] Building API bridge...
go build -o api-bridge-test.exe main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build API bridge
    pause
    exit /b 1
)

echo [INFO] API bridge built successfully

REM Start API bridge in background
echo [INFO] Starting API bridge with event system enabled...
start /B api-bridge-test.exe -config test_config.yaml

REM Wait for API bridge to start
echo [INFO] Waiting for API bridge to start...
timeout /t 5 /nobreak >nul

REM Test if API bridge is running
curl -s http://localhost:8081/health >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Failed to start API bridge
    echo Please check if port 8081 is available
    pause
    exit /b 1
)

echo [INFO] API bridge started successfully

echo.
echo [INFO] Testing event system...
echo.

REM Test component registration with event emission
echo [INFO] Testing component registration with event emission...

curl -s -X POST http://localhost:8081/api/v1/components/register ^
    -H "Content-Type: application/json" ^
    -d "{\"creator\": \"test-user\", \"component_data\": \"test-battery-module\", \"context\": \"event-system-test\"}"

if %errorlevel% equ 0 (
    echo [INFO] Component registration successful
) else (
    echo [ERROR] Component registration failed
)

REM Wait a moment for event processing
timeout /t 3 /nobreak >nul

REM Test pairing initiation with event emission
echo [INFO] Testing pairing initiation with event emission...

curl -s -X POST http://localhost:8081/api/v1/pairing/initiate ^
    -H "Content-Type: application/json" ^
    -d "{\"creator\": \"test-user\", \"component_a\": \"COMP-test-user-123\", \"component_b\": \"COMP-test-user-456\", \"operational_context\": \"event-system-test\", \"force_immediate\": true}"

if %errorlevel% equ 0 (
    echo [INFO] Pairing initiation successful
) else (
    echo [ERROR] Pairing initiation failed
)

REM Wait a moment for event processing
timeout /t 3 /nobreak >nul

echo.
echo [INFO] Checking for received events...
echo.
echo === Expected Events ===
echo The following events should have been sent to webhook endpoints:
echo.
echo 1. component_registered event:
echo    - Endpoint: http://localhost:3000/webhooks/component-registered
echo    - Contains: component_id, creator, component_data, context, timestamp, tx_hash
echo.
echo 2. pairing_initiated event:
echo    - Endpoint: http://localhost:3000/webhooks/pairing-initiated
echo    - Contains: challenge_id, creator, component_a, component_b, operational_context, timestamp, tx_hash
echo.
echo Note: In a real deployment, you would:
echo - Set up proper webhook receivers
echo - Store events in SQL databases
echo - Implement authentication and security
echo - Add monitoring and alerting
echo.
echo [INFO] Event system test completed!
echo [INFO] Check the API bridge logs for event emission details
echo.
echo [WARN] Press any key to stop the test and cleanup
pause

REM Cleanup
echo [INFO] Cleaning up...

REM Stop API bridge
taskkill /F /IM api-bridge-test.exe >nul 2>&1

REM Clean up test files
del api-bridge-test.exe >nul 2>&1
del test_config.yaml >nul 2>&1

echo [INFO] Cleanup completed
pause 