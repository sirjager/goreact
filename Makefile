# grpc and gateways
SERVICE_NAME=goreact
STUBS_DIR=./stubs
PROTO_DIR=./proto

GO_STUB_DIR=$(STUBS_DIR)/go
JS_STUB_DIR=$(STUBS_DIR)/js

STATIK_DIR=./statik
# docs
DOCS_PKG=docs
DOCS_SWAGGER=$(STATIK_DIR)/swagger
# web app
WEB_PKG="web"
WEB_SRC=./app
WEB_OUT=$(STATIK_DIR)/$(WEB_PKG)


# TEST: database configs 
DB_MIGRATIONS=./migrations
DB_URL=postgres://postgres:password@localhost:5432/$(SERVICE_NAME)?sslmode=disable

# Generate GO protobuf and gRPC code, Swagger documentation, and static assets
proto:
	- rm -rf $(GO_STUB_DIR)
	- mkdir -p $(GO_STUB_DIR)
	- rm -f $(DOCS_SWAGGER)/*.swagger.json
	- rm -rf $(STATIK_DIR)/$(DOCS_PKG)
	protoc \
	--proto_path=$(PROTO_DIR) --go_out=$(GO_STUB_DIR) --go_opt=paths=source_relative	\
	--go-grpc_out=$(GO_STUB_DIR) --go-grpc_opt=paths=source_relative	\
	--grpc-gateway_out=$(GO_STUB_DIR) --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=$(DOCS_SWAGGER) --openapiv2_opt=allow_merge=true,merge_file_name=$(SERVICE_NAME) \
	$(PROTO_DIR)/*.proto
	statik -src=$(DOCS_SWAGGER) -dest=$(STATIK_DIR) -ns=$(DOCS_PKG) -p=$(DOCS_PKG)

web:
	cd $(WEB_SRC) && bun run build
	- rm -rf $(STATIK_DIR)/$(WEB_PKG)
	statik -src=$(WEB_SRC)/dist -dest=$(STATIK_DIR) -ns=$(WEB_PKG) -p=$(WEB_PKG)

proto-js:
	- rm -rf $(JS_STUB_DIR)
	- mkdir -p $(JS_STUB_DIR) 
	pnpm grpc_tools_node_protoc	\
	--js_out=import_style=commonjs,binary:${JS_STUB_DIR} \
	--grpc_out=${JS_STUB_DIR} \
	--plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
	-I $(PROTO_DIR)/ \
	$(PROTO_DIR)/*.proto
	
	# Generate typescript (d.ts)
	protoc \
	--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:$(JS_STUB_DIR) \
	-I $(PROTO_DIR) \
	$(PROTO_DIR)/*.proto


# Clean up go.sum and tidy dependencies
tidy:
	- rm -f ./go.sum
	- go mod tidy

# Run tests
test:
	go clean -testcache
	go test -v -cover -short ./... 

# Build the project
build:
	golint ./...
	go build -o ./dist/main ./cmd/main.go

# Linting and formatting
lint:
	golint ./...

# Run the project
run:
	go run ./cmd/main.go

# Run the project
dev:
	air && fg

bindata:
	cd ./migrations && go-bindata -pkg migrations .

# Generate database documentation
dbdocs:
	dbdocs build ./internal/db/db.dbml;

# Generate database schema
dbschema:
	dbml2sql --postgres -o ./internal/db/schema.sql ./internal/db/db.dbml

# Create a new database migration
migratenew:
	migrate create -ext sql -dir $(DB_MIGRATIONS) -seq $(filter-out $@,$(MAKECMDGOALS))

# Run database migrations (up)
migrateup:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose up

# Rollback database migrations (down)
migratedown:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose down --all

# Generate SQL code using SQLC
sqlc:
	sqlc generate

release:
	goreleaser --snapshot --clean

image:
	docker build --pull --rm -f "Dockerfile" -t gomicro:latest "."

deps:
	go install github.com/cosmtrek/air@latest

# Declare phony targets to prevent conflicts with file names
.PHONY: proto proto-js tidy test run dbdocs dbschema migratenew migrateup migratedown sqlc release
