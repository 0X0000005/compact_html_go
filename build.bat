@echo off
REM compact_html_go - Cross-platform build script

SET VERSION=1.0.0
SET APP_NAME=compact_html

echo ================================================
echo compact_html_go - Build script
echo Version: %VERSION%
echo ================================================
echo.

echo [1/5] Cleaning old files...
if exist %APP_NAME%.exe del %APP_NAME%.exe
if exist %APP_NAME% del %APP_NAME%

echo.
echo [2/5] Preparing dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Dep download failed
    exit /b 1
)

echo.
echo [3/5] Compiling Windows version...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o %APP_NAME%.exe .\cmd\main.go
if errorlevel 1 (
    echo Error: Windows compilation failed
    exit /b 1
)
echo OK

echo.
echo [4/5] Compiling Linux version...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w -X main.version=%VERSION%" -o %APP_NAME% .\cmd\main.go
set GOOS=
set GOARCH=
if errorlevel 1 (
    echo Error: Linux compilation failed
    exit /b 1
)
echo OK

echo.
echo [5/5] UPX compression...
where upx >nul 2>nul
if errorlevel 1 (
    echo Warning: UPX not found, skipping
    goto :finish
)

echo Compressing Windows version...
upx -9 %APP_NAME%.exe
if errorlevel 1 (
    echo Warning: Windows version compression failed
)

echo Compressing Linux version...
upx -9 %APP_NAME%
if errorlevel 1 (
    echo Warning: Linux version compression failed
)

:finish
echo.
echo ================================================
echo Build finished!
echo ================================================
echo Windows: %APP_NAME%.exe
echo Linux:   %APP_NAME%
echo.
echo Start command (Windows): %APP_NAME%.exe -h
echo ================================================
echo.
