@echo off
chcp 65001
cls

:: 保存当前目录并切换到脚本所在目录
pushd "%~dp0"

:: 调用主逻辑子例程
call :main %1 %2 %3

:: 恢复原始目录
popd

:: 退出并返回错误码
exit /b %errorlevel%

:main

set app_name=jvms

:: 清git理残留的编译进程
taskkill /f /im go.exe    >nul 2>nul
taskkill /f /im compile.exe    >nul 2>nul
taskkill /f /im asm.exe        >nul 2>nul
taskkill /f /im link.exe       >nul 2>nul
taskkill /f /im git.exe        >nul 2>nul
taskkill /f /im %app_name%.exe >nul 2>nul

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

echo =============================================================
echo 1 - 编译 Windows 可执行程序
echo =============================================================
go mod tidy
SET GO111MODULE=on
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o %app_name%.exe -ldflags "-X main.IsBeta=%IsBeta% -X main.BuildTime=%BuildTime%" .
if errorlevel 1 (
    echo 编译失败，请检查错误信息。
    exit /b 1
)

echo 2 - 运行程序
echo =============================================================
echo %app_name% rls
