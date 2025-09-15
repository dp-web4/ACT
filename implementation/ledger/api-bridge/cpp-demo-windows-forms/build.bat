@echo off
setlocal enabledelayedexpansion

echo ========================================
echo Web4 API Bridge Demo - Windows Forms
echo Building for Windows...
echo ========================================

:: Check if we're in the right directory
if not exist "CMakeLists.txt" (
    echo Error: CMakeLists.txt not found. Please run this script from the project root directory.
    pause
    exit /b 1
)

:: Set build directory
set BUILD_DIR=build
set BIN_DIR=%BUILD_DIR%\bin

:: Create build directory if it doesn't exist
if not exist "%BUILD_DIR%" (
    echo Creating build directory...
    mkdir "%BUILD_DIR%"
)

:: Check for Visual Studio
set VS_FOUND=0
set VS_VERSION=""

:: Check for Visual Studio 2022
where "C:\Program Files\Microsoft Visual Studio\2022\*\MSBuild\Current\Bin\MSBuild.exe" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    set VS_FOUND=1
    set VS_VERSION="Visual Studio 17 2022"
    echo Found Visual Studio 2022
    goto :build
)

:: Check for Visual Studio 2019
where "C:\Program Files (x86)\Microsoft Visual Studio\2019\*\MSBuild\Current\Bin\MSBuild.exe" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    set VS_FOUND=1
    set VS_VERSION="Visual Studio 16 2019"
    echo Found Visual Studio 2019
    goto :build
)

:: Check for Visual Studio 2017
where "C:\Program Files (x86)\Microsoft Visual Studio\2017\*\MSBuild\15.0\Bin\MSBuild.exe" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    set VS_FOUND=1
    set VS_VERSION="Visual Studio 15 2017"
    echo Found Visual Studio 2017
    goto :build
)

:: No Visual Studio found
echo Warning: No Visual Studio found. Will use default generator.
set VS_FOUND=0
set VS_VERSION=""

:build

:: Check for vcpkg
set VCPKG_FOUND=0
if exist "C:\vcpkg\scripts\buildsystems\vcpkg.cmake" (
    set VCPKG_FOUND=1
    set VCPKG_TOOLCHAIN="-DCMAKE_TOOLCHAIN_FILE=C:/vcpkg/scripts/buildsystems/vcpkg.cmake"
    echo Found vcpkg at C:\vcpkg
) else if exist "%USERPROFILE%\vcpkg\scripts\buildsystems\vcpkg.cmake" (
    set VCPKG_FOUND=1
    set VCPKG_TOOLCHAIN="-DCMAKE_TOOLCHAIN_FILE=%USERPROFILE%/vcpkg/scripts/buildsystems/vcpkg.cmake"
    echo Found vcpkg at %USERPROFILE%\vcpkg
) else (
    echo Warning: vcpkg not found. Some dependencies may not be available.
    set VCPKG_TOOLCHAIN=""
)

:: Configure with CMake
echo.
echo Configuring with CMake...
cd "%BUILD_DIR%"

if %VS_FOUND% EQU 1 (
    if %VCPKG_FOUND% EQU 1 (
        cmake .. -G %VS_VERSION% -A x64 %VCPKG_TOOLCHAIN%
    ) else (
        cmake .. -G %VS_VERSION% -A x64
    )
) else (
    if %VCPKG_FOUND% EQU 1 (
        cmake .. %VCPKG_TOOLCHAIN%
    ) else (
        cmake ..
    )
)

if %ERRORLEVEL% NEQ 0 (
    echo Error: CMake configuration failed.
    cd ..
    pause
    exit /b 1
)

:: Build the project
echo.
echo Building project...
cmake --build . --config Release

if %ERRORLEVEL% NEQ 0 (
    echo Error: Build failed.
    cd ..
    pause
    exit /b 1
)

:: Check if executable was created
cd ..
if not exist "%BIN_DIR%\APIBridgeDemoWindowsForms.exe" (
    echo Error: Executable not found at %BIN_DIR%\APIBridgeDemoWindowsForms.exe
    pause
    exit /b 1
)

echo.
echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo Executable location: %BIN_DIR%\APIBridgeDemoWindowsForms.exe
echo.
echo To run the application:
echo   cd %BIN_DIR%
echo   APIBridgeDemoWindowsForms.exe
echo.

:: Ask if user wants to run the application
set /p RUN_APP="Do you want to run the application now? (y/n): "
if /i "%RUN_APP%"=="y" (
    echo.
    echo Starting application...
    cd "%BIN_DIR%"
    start APIBridgeDemoWindowsForms.exe
    cd ..
)

echo.
echo Build script completed.
pause 