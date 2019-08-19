@echo off

echo Compiling Protocols
echo NOTE: Make sure you have 'protoc.exe' extracted to your path as well
echo          as having 'github.com/golang/protobuf' installed in go.
echo (1) choco install protoc
echo (2) go get -d -u github.com/golang/protobuf/protoc-gen-go
echo (3) go install github.com/golang/protobuf/protoc-gen-go
echo (4) Also, ensure that your GO/bin is added to path

pushd %cd%
cd pb/protos

for /F "" %%i in ('dir /b *.proto') do protoc --go_out=../ %%i
popd
echo Done!