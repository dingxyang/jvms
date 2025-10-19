@echo off
chcp 65001
cls

set app_name=jvms
set app_ver=3.0.6

rem 获取当前时间
set "hour=%time:~0,2%"
if "%hour:~0,1%" == " " set "hour=0%hour:~1,1%"

set "year=%date:~0,4%"
set "month=%date:~5,2%"
set "day=%date:~8,2%"

set "minute=%time:~3,2%"
set "second=%time:~6,2%"

set "date_text=%year%-%month%-%day%(%hour%:%minute%:%second%)"
set "date_version=%year%%month%%day%_%hour%%minute%%second%

:: 发布正式版本时，去掉版本号中的时间信息
set BuildTime=%date_text%
echo 编译时间：%date_text%

:: 清理残留的编译进程
taskkill /f /im go.exe         >nul 2>nul
taskkill /f /im compile.exe    >nul 2>nul
taskkill /f /im asm.exe        >nul 2>nul
taskkill /f /im link.exe       >nul 2>nul
taskkill /f /im git.exe        >nul 2>nul

del %app_name%*.zip                  >nul 2>nul
del %app_name%*.tar.gz               >nul 2>nul

:: 编译程序
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
taskkill /f /im %app_name%.exe   >nul 2>nul
del %app_name%.exe               >nul 2>nul
attrib -H *.old                  >nul 2>nul
del *.exe.old                    >nul 2>nul

go build -o %app_name%.exe -ldflags "-X main.AppVersion=%app_ver% -X main.BuildTime=%BuildTime%" .
if errorlevel 1 (
    echo 编译失败，请检查错误信息。
    exit /b 1
)

:: 获取版本号信息
:: echo %app_name%.exe -v
tar -a -c -f %app_name%%app_ver%.windows-amd64.zip %app_name%.exe >nul 2>nul
%app_name%.exe -v >nul 2>nul
if errorlevel 1 (
    echo [E] 获取版本失败，请确认是否有 -v 参数。
    exit /b 1
)

for /f "delims=" %%a in ('%app_name%.exe -v') do @set "app_version=%%a"
echo 当前版本：%app_version%
echo =============================================================

echo 1 - 编译 Windows 可执行程序
echo =============================================================

if "%1" == "test" (
    echo %app_name%.exe -v
    %app_name%.exe -v
    exit /b 1
)

echo 2 - 编译 Linux 可执行程序
echo =============================================================
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
del %app_name% >nul 2>nul

go build -o %app_name% -ldflags "-X main.AppVersion=%app_ver% -X main.BuildTime=%BuildTime%"  .
tar -czvf %app_name%%app_ver%.linux-amd64.tar.gz %app_name% >nul 2>nul

echo 3 - 编译 OpenWRT 可执行程序
echo =============================================================
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
del %app_name% >nul 2>nul

go build -o %app_name% -ldflags "-X main.AppVersion=%app_ver% -X main.BuildTime=%BuildTime%"  .
tar -czvf %app_name%%app_ver%.linux-arm64.tar.gz %app_name% >nul 2>nul

echo 3 - 编译 MacOS (ARM64) 可执行程序
echo =============================================================
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=arm64
del %app_name% >nul 2>nul

go build -o %app_name% -ldflags "-X main.AppVersion=%app_ver% -X main.BuildTime=%BuildTime%"  .
tar -czvf %app_name%%app_ver%.darwin-arm64.tar.gz %app_name% >nul 2>nul

echo 3 - 编译 MacOS (X86) 可执行程序
echo =============================================================
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
del %app_name% >nul 2>nul

go build -o %app_name% -ldflags "-X main.AppVersion=%app_ver% -X main.BuildTime=%BuildTime%"  .
tar -czvf %app_name%%app_ver%.darwin-amd64.tar.gz %app_name% >nul 2>nul
del %app_name% >nul 2>nul