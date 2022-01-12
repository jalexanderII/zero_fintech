.PHONY: gen_core clear_core migration_core migrate_core rollback_core drop_core

gen_core:
	protoc -I=./proto/core --go_opt=paths=source_relative --go_out=plugins=grpc:./gen/core ./proto/core/*.proto

clear_core:
	rm gen/core/*.go

migration_core:
	@read -p "Enter migration name: " name; migrate create -ext sql -dir database/migrations $$name

migrate_core:
	migrate -source=file://database/migrations -database postgres://postgres:Qnacmg797y@127.0.0.1:5432/coredb?sslmode=disable up

rollback_core:
	migrate -source=file://database/migrations -database postgres://postgres:Qnacmg797y@127.0.0.1:5432/coredb?sslmode=disable down

drop_core:
	migrate -source=file://database/migrations -database postgres://postgres:Qnacmg797y@127.0.0.1:5432/coredb?sslmode=disable drop