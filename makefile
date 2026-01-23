include .env

LOCAL_BIN = $(CURDIR)/shared/bin

# ==========================================================================================================================

install-all:
	make install-grpc
	make install-protoc
	make install-grpc-gateway
	make install-mockery
	make install-goose

install-grpc:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10

install-protoc:
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.0

install-grpc-gateway:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.3

install-wire:
	GOBIN=$(LOCAL_BIN) go install github.com/google/wire/cmd/wire@v0.7.0

# ==========================================================================================================================

get-annotation:
	git clone https://github.com/googleapis/googleapis.git shared/vdr/googletmp &&\
	mkdir -p shared/vdr &&\
	mv shared/vdr/googletmp/google shared/vdr &&\
	rm -rf shared/vdr/googletmp;\

# ==========================================================================================================================

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

generate-analytics:
	mkdir -p shared/gen/analytics
	protoc --proto_path api --proto_path shared/vdr \
	--go_out=shared/gen/analytics --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=shared/bin/protoc-gen-go \
	--go-grpc_out=shared/gen/analytics --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=shared/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=shared/gen/analytics --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=shared/bin/protoc-gen-grpc-gateway \
	api/analytics.proto


# ==========================================================================================================================

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

# ==========================================================================================================================

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


migration-create-analytics:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="${PG_LINKS_DSN}" \
	${LOCAL_BIN}/goose -dir ${ANALYTICS_MIGRATION_DIR} create /analytics_table sql

migrations-up-analytics:
	${LOCAL_BIN}/goose -dir ${ANALYTICS_MIGRATION_DIR} postgres ${PG_LINKS_DSN} -table goose_version_analytics up -v

migrations-down-analytics:
	${LOCAL_BIN}/goose -dir ${ANALYTICS_MIGRATION_DIR} postgres ${PG_ANALYTICS_DSN} -table goose_version_analytics down -v

# ==========================================================================================================================

install-mockery:
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v3@v3.6.1
	
generate-mocks:
	${LOCAL_BIN}/mockery --config mockery_auth.yaml
	${LOCAL_BIN}/mockery --config mockery_links.yaml
	${LOCAL_BIN}/mockery --config mockery_shared.yaml

delete-mocks:
	rm -r ./shared/mocks
	rm -r ./services/auth/internal/mocks
	rm -r ./services/links/internal/mocks

# ==========================================================================================================================

tests:
	docker compose -f docker-compose_test.yaml up -d
	sleep 2
	go clean -testcache
	cd ./services/auth && go test ./...
	cd ./services/links && go test ./...
	docker stop url_shortener_test
	docker rm url_shortener_test
	

nats-box:
	docker run -it --rm --network url_shortener_default natsio/nats-box:latest sh
