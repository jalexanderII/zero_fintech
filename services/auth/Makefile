.PHONY: gen_auth clear_auth

gen_auth:
	protoc -I=./proto/auth --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/auth ./proto/auth/*.proto

clear_auth:
	rm gen/auth/*.go
