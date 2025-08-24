include .env

.PHONY: generate
generate: sqlc swagger mock

.PHONY: sqlc
sqlc:
	sqlc generate
	find ./generate/*sql -type f \( -name "models.go" -o -name "*.sql.go" -o -name "querier.go" \) -exec sed -i '' \
		-e 's/\"database\/sql\"//g' \
		-e 's/string/\*string/g' \
		-e 's/\([^:]\)bool/\1*bool/g' \
		-e 's/int64/\*uint64/g' \
		-e 's/int32/\*int32/g' \
		-e 's/float64/\*float64/g' \
		-e 's/time\.Time/\*time.Time/g' \
		-e 's/sql\.NullString/\*string/g' \
		-e 's/sql\.NullTime/\*time.Time/g' \
		-e 's/sql\.NullInt64/\*uint64/g' \
		-e 's/sql\.NullInt32/\*int32/g' \
		-e 's/sql\.NullFloat64/\*float64/g' \
		-e 's/sql\.NullBool/\*bool/g' \
		{} +
	find ./generate/psql -type f -exec sed -i '' -e 's/DBTX/PDBTX/g' -e 's/Querier/PQuerier/g' {} +
	pg_dump --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels --exclude-table=_gooses ${GOOSE_DBSTRING} > generate/raw.sql
	awk '/^CREATE TABLE/,/^\);$$/' generate/raw.sql > generate/schema.sql

.PHONY: swagger
swagger:
	go run ./command/swagger
	swag init -g ./common/swagger/swagger.go -o ./generate/swagger/ --parseDependency
	sed -i '' -e 's/response.GenericResponse/Res/g' -e 's/payload_//g' ./generate/swagger/*.{json,yaml}
	swagger2openapi --yaml --outfile ./generate/swagger/openapi.yaml ./generate/swagger/swagger.yaml
	go run ./command/swaggermd
	(cd ../frontend && bun run generate)

.PHONY: mock
mock:
	mockery

.PHONY: build
build:
	make generate
	nilaway ./... || true
	(cd ../frontend && bun run build)
	cp -r ../frontend/dist/ .local/dist/
	env GOOS=linux GOARCH=amd64 go build -o .local/backend .