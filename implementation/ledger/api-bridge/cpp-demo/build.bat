@echo off
setlocal enabledelayedexpansion

echo ========================================
echo Web4 API Bridge Demo - Build Script
echo ========================================

:: Check if CMake is available
where cmake >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: CMake not found. Please install CMake and add it to PATH.
    pause
    exit /b 1
)

:: Check if Git is available
where git >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Git not found. Please install Git and add it to PATH.
    pause
    exit /b 1
)

:: Set build directory
set BUILD_DIR=build
set VCPKG_DIR=C:\vcpkg

:: Check if vcpkg exists
if not exist "%VCPKG_DIR%" (
    echo Installing vcpkg...
    git clone https://github.com/Microsoft/vcpkg.git "%VCPKG_DIR%"
    cd "%VCPKG_DIR%"
    call bootstrap-vcpkg.bat
    call vcpkg integrate install
    cd ..
)

:: Install dependencies
echo Installing dependencies...
cd "%VCPKG_DIR%"
call vcpkg install nlohmann-json:x64-windows
call vcpkg install grpc:x64-windows
call vcpkg install protobuf:x64-windows
cd ..

:: Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
cd "%BUILD_DIR%"

:: Configure with CMake
echo Configuring with CMake...
cmake .. -DCMAKE_TOOLCHAIN_FILE="%VCPKG_DIR%\scripts\buildsystems\vcpkg.cmake" -DCMAKE_BUILD_TYPE=Release
if %errorlevel% neq 0 (
    echo ERROR: CMake configuration failed.
    pause
    exit /b 1
)

:: Build the project
echo Building project...
cmake --build . --config Release
if %errorlevel% neq 0 (
    echo ERROR: Build failed.
    pause
    exit /b 1
)

:: Copy configuration file
if not exist "bin\config.json" (
    copy "..\config.json" "bin\"
)

echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo Executable location: %BUILD_DIR%\bin\APIBridgeDemo.exe
echo.
echo To run the demo:
echo   cd %BUILD_DIR%\bin
echo   APIBridgeDemo.exe
echo.
echo To run with custom config:
echo   APIBridgeDemo.exe --config config.json
echo.

:: Ask if user wants to run the demo
set /p RUN_DEMO="Do you want to run the demo now? (y/n): "
if /i "%RUN_DEMO%"=="y" (
    echo Starting demo...
    cd bin
    APIBridgeDemo.exe
)

pause 