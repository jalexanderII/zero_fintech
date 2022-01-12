Run

python3 -m grpc_tools.protoc -I proto/planning --python_out=. --grpc_python_out=. proto/planning/*.proto