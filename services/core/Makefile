.PHONY: gen_core clear_core

gen_core:
	protoc -I=./proto/core --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/core ./proto/core/*.proto

clear_core:
	rm gen/core/*.go
