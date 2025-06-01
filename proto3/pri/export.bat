@echo off
for %%i in (*.proto) do (
    protoc ./%%i --gogo_out=../
    echo exchange %%i To go file successfully!
)