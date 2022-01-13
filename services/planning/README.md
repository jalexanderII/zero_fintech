Run

<!-- python3 -m grpc_tools.protoc -I proto/planning --python_out=. --grpc_python_out=. proto/planning/*.proto -->

python3 -m grpc_tools.protoc -I proto/planning --python_out=gen/planning --grpc_python_out=gen/planning proto/planning/*.proto