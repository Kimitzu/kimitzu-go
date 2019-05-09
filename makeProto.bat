@echo off

echo Compiling Protocols
echo NOTE: Make sure you have 'protoc.exe' extracted to your path as well
echo          as having 'github.com/golang/protobuf' installed in go.

pushd %cd%
cd pb/protos

for /F "" %%i in ('dir /b *.proto') do protoc --go_out=../ %%i
popd
echo Done!