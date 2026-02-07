@echo off

:: Cek apakah ada argumen yang dikirim
if "%1"=="" goto help

:: Pindah ke label sesuai argumen pertama
if "%1"=="gen" goto gen

:gen
    if "%2"=="" (echo Error: Perlu nama migrasi! Contoh: setup gen nama_file & exit /b)
    docker run --rm -v "%cd%/db/migrations:/migrations" migrate/migrate:v4.19.1 create -ext sql -dir /migrations -seq %2
    goto end
