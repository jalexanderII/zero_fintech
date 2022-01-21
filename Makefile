.PHONY: gen_protos clear_protos

gen_protos_go:
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/common/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/auth/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/planning/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/core/*.proto

clear_protos_go:
	rm ./gen/Go/common/*.go
	rm ./gen/Go/auth/*.go
	rm ./gen/Go/planning/*.go
	rm ./gen/Go/core/*.go

gen_protos_py:
	python3 -m grpc_tools.protoc -I proto --python_betterproto_out=services/planning/database/models --python_out=gen/Python --mypy_out=gen/Python --grpc_python_out=gen/Python proto/auth/*.proto
	python3 -m grpc_tools.protoc -I proto --python_betterproto_out=services/planning/database/models --python_out=gen/Python --mypy_out=gen/Python --grpc_python_out=gen/Python proto/common/*.proto
	python3 -m grpc_tools.protoc -I proto --python_betterproto_out=services/planning/database/models --python_out=gen/Python --mypy_out=gen/Python --grpc_python_out=gen/Python proto/core/*.proto
	python3 -m grpc_tools.protoc -I proto --python_betterproto_out=services/planning/database/models --python_out=gen/Python --mypy_out=gen/Python --grpc_python_out=gen/Python proto/planning/*.proto

clear_protos_py:
	rm ./gen/Python/common/*.py
	rm ./gen/Python/auth/*.py
	rm ./gen/Python/planning/*.py
	rm ./gen/Python/core/*.py
