@echo off

goto %1
goto end

:clrp
echo # This will delete "%userprofile%\openbazaar" prior to starting Djali.
set /p conf="# Continue?[y/N]: "

if /i %conf%==y goto delete
goto end

:delete
rm -fdr "%userprofile%\openbazaar"
echo # Wiped

:start
go run openbazaard.go start --verbose
goto :end

:hello
echo hello

:end
echo # End