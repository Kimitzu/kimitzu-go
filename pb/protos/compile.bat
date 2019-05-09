@echo off

echo Compiling Protocols
echo NOTE: Make sure you're running this in the 'protos' directory
echo          and that you have 'protoc.exe' extracted to your path as well
echo          as having 'github.com/golang/protobuf' installed in go.

for /F "" %%i in ('dir /b *.proto') do protoc --go_out=../ %%i

echo Done!