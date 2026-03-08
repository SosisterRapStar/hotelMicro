-include .env
.EXPORT_ALL_VARIABLES:

.PHONY: install-migrate
install-migrate:
	go install -tags 'mysql,postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migration-create
migration-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migration-up
migration-up:
	migrate -path migrations -database "$(MIGRATE_DSN)" up

.PHONY: migration-down
migration-down:
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

.PHONY: migration-force
migration-force:
	migrate -path migrations -database "$(MIGRATE_DSN)" force $(version)

.PHONY: run
run:
	go run ./cmd/app.go

.PHONY: proto
proto:
	@buf generate

.PHONY: swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.8.1 init -g cmd/app.go -o internal/docs --parseDependency --parseInternal