@echo on
set BINDIR=%CD%
pushd "C:\Program Files\SAP Lumira\Desktop\daextensions"
ls
del %1.exe
xcopy %BINDIR%\%1.exe
ls
