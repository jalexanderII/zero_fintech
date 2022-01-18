.PHONY: gen_protos clear_protos

gen_protos:
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/common/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/auth/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/planning/*.proto
	protoc -I=./proto --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/Go/ ./proto/core/*.proto

clear_protos:
	rm ./gen/Go/common/*.go
	rm ./gen/Go/auth/*.go
	rm ./gen/Go/planning/*.go
	rm ./gen/Go/core/*.go
