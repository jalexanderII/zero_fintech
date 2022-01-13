# Zero FinTech Core Service

python3 -m grpc_tools.protoc -I proto/core --python_out=gen/core --grpc_python_out=gen/core proto/core/*.proto


python3 -m grpc_tools.protoc -I proto/core --python_out=gen/core_py --grpc_python_out=gen/core_py proto/core/*.proto