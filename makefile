include .env

LOCAL_BIN = $(CURDIR)/shared/bin

# ==========================================================================================================================

install-all:
	make install-grpc
	make install-protoc
	make install-grpc-gateway
	make install-goose
	make install-mockery

install-grpc:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10

install-protoc:
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.0

install-grpc-gateway:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.3

install-wire:
	GOBIN=$(LOCAL_BIN) go install github.com/google/wire/cmd/wire@v0.7.0

install-mockery:
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v3@v3.6.1

# ==========================================================================================================================

get-annotation:
	git clone https://github.com/googleapis/googleapis.git shared/vdr/googletmp &&\
	mkdir -p shared/vdr &&\
	mv shared/vdr/googletmp/google shared/vdr &&\
	rm -rf shared/vdr/googletmp;\

# ==========================================================================================================================

generate-services:
	make generate-auth
	make generate-links

generate-mocks:
	${LOCAL_BIN}/mockery --config .mockery.yml

generate-auth:
	mkdir -p shared/gen/auth
	protoc --proto_path api --proto_path shared/vdr \
	--go_out=shared/gen/auth --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=shared/bin/protoc-gen-go \
	--go-grpc_out=shared/gen/auth --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=shared/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=shared/gen/auth --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=shared/bin/protoc-gen-grpc-gateway \
	api/auth.proto

generate-links:
	mkdir -p shared/gen/links
	protoc --proto_path api --proto_path shared/vdr \
	--go_out=shared/gen/links --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=shared/bin/protoc-gen-go \
	--go-grpc_out=shared/gen/links --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=shared/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=shared/gen/links --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=shared/bin/protoc-gen-grpc-gateway \
	api/links.proto

# ==========================================================================================================================

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

# ==========================================================================================================================

migrations-up:
	make migrations-up-auth
	make migrations-up-links

migrations-down:
	make migrations-down-auth
	make migrations-down-links

migration-create-auth:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="${PG_AUTH_DSN}" \
	${LOCAL_BIN}/goose -dir ${AUTH_MIGRATION_DIR} create auth_table sql

migrations-up-auth:
	${LOCAL_BIN}/goose -dir ${AUTH_MIGRATION_DIR} postgres ${PG_AUTH_DSN} -table goose_version_auth up -v

migrations-down-auth:
	${LOCAL_BIN}/goose -dir ${AUTH_MIGRATION_DIR} postgres ${PG_AUTH_DSN} -table goose_version_auth down -v


migration-create-links:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="${PG_LINKS_DSN}" \
	${LOCAL_BIN}/goose -dir ${LINKS_MIGRATION_DIR} create links_table sql

migrations-up-links:
	${LOCAL_BIN}/goose -dir ${LINKS_MIGRATION_DIR} postgres ${PG_LINKS_DSN} -table goose_version_links up -v

migrations-down-links:
	${LOCAL_BIN}/goose -dir ${LINKS_MIGRATION_DIR} postgres ${PG_LINKS_DSN} -table goose_version_links down -v


app-setup:
	make install-all
	make get-annotation
	make generate-services
	make generate-mocks
	clear
	@echo "\033[32m✅ Setup done\033[0m"
