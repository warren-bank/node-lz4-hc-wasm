@echo off

set CWD=%CD%
set PRJ_DIR=%~dp0..\..

if not exist "%PRJ_DIR%\dist" mkdir "%PRJ_DIR%\dist"
if exist "%PRJ_DIR%\dist\lz4.wasm" del "%PRJ_DIR%\dist\lz4.wasm"
if exist "%PRJ_DIR%\dist\wasm_exec.js" del "%PRJ_DIR%\dist\wasm_exec.js"

rem :: ---------------------------------------------------------------

set GOOS=js
set GOARCH=wasm
set GO111MODULE=on

cd /D "%PRJ_DIR%\src"
call go build -o "%PRJ_DIR%\dist\lz4.wasm" "main.go"

rem :: ---------------------------------------------------------------

copy "%GOROOT%\lib\wasm\wasm_exec.js" "%PRJ_DIR%\dist"

rem :: ---------------------------------------------------------------

cd /D "%CWD%"
