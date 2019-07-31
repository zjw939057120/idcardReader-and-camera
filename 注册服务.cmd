@echo off
set filepath=%~dp0
set file=bin\go_build_server_go.exe
echo server:  %filepath%%file%
echo step:  1
reg add "HKEY_CLASSES_ROOT\Reader" /ve /t REG_SZ /d "reader Protocol" /f
echo step:  2
reg add "HKEY_CLASSES_ROOT\Reader" /v "URL Protocol" /t REG_SZ /f
echo step:  3
reg add "HKEY_CLASSES_ROOT\Reader\DefaultIcon" /ve /t REG_EXPAND_SZ /d %filepath%%file% /f
echo step:  4
reg add "HKEY_CLASSES_ROOT\Reader\shell\open\command" /ve /t REG_EXPAND_SZ /d %filepath%%file% /f


cmd /k echo Done
