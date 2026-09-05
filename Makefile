dbUrl ?= $(or $(MIGRATE_DATABASE_URL),postgres://myuser:mypassword@localhost:5433/myapp?sslmode=disable)
composeFile = "./deploy/docker-compose.yml"
dockerFile = "deploy/Dockerfile"
atlasDbUrl ?= $(or $(ATLAS_DATABASE_URL),postgres://myuser:mypassword@postgres:5432/myapp?sslmode=disable)
atlasDevUrl ?= $(or $(ATLAS_DEV_DATABASE_URL),postgres://myuser:mypassword@postgres:5432/atlas_dev?sslmode=disable)

dev: 
	@docker-compose --env-file .env -f $(composeFile) up -d

down: 
	@docker-compose --env-file .env -f $(composeFile) down

logs: 
	@docker-compose --env-file .env -f $(composeFile) logs -f 

migration:
	@migrate create -ext sql -dir internal/db/migrations -seq $(name)

migrate:
	@migrate -database '$(dbUrl)' -path internal/db/migrations up

rollback:
	@migrate -database '$(dbUrl)' -path internal/db/migrations down 1

clean-db:
	@set -e; \
	version_output="$$(migrate -database '$(dbUrl)' -path internal/db/migrations version 2>&1 || true)"; \
	case "$$version_output" in \
		*"(dirty)"*) \
			dirty_version="$${version_output%% *}"; \
			echo "database is dirty at version $$dirty_version; forcing clean migration state"; \
			migrate -database '$(dbUrl)' -path internal/db/migrations force $$dirty_version; \
			;; \
	esac; \
	migrate -database '$(dbUrl)' -path internal/db/migrations down -all

reset-db: clean-db migrate

test:
	@go test ./...

gen-query:
	@go run ./cmd/querygen

verify-gen-query:
	@go run ./cmd/querygen
	@git diff --exit-code -- internal/db/query cmd/querygen/main.go internal/db/persistence/schema.go

atlas-version:
	@ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) run --rm atlas version

atlas-reset-dev-db:
	@ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) exec -T postgres sh -lc 'psql -U "$$POSTGRES_USER" -d postgres -c "DROP DATABASE IF EXISTS atlas_dev;"'
	@ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) exec -T postgres sh -lc 'psql -U "$$POSTGRES_USER" -d postgres -c "CREATE DATABASE atlas_dev;"'

atlas-schema:
	@$(MAKE) atlas-reset-dev-db
	@ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) run --rm atlas schema inspect --env local --url env://src --format '{{ sql . }}'

atlas-diff:
	@test -n "$(name)" || (echo "usage: make atlas-diff name=<migration_name>" && exit 1)
	@set -e; \
	$(MAKE) atlas-reset-dev-db; \
	diff_sql="$$(ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) run --rm atlas schema diff --env local --from "$(atlasDbUrl)" --to env://src --exclude "public.schema_migrations" --format '{{ sql . }}')"; \
	if [ -z "$$diff_sql" ]; then \
		echo "no schema changes detected; no migration files created"; \
		exit 0; \
	fi; \
	reverse_sql="$$(ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) run --rm atlas schema diff --env local --from env://src --to "$(atlasDbUrl)" --exclude "public.schema_migrations" --format '{{ sql . }}')"; \
	migrate create -ext sql -dir internal/db/migrations -seq $(name); \
	up_file=$$(ls -1 internal/db/migrations/*_$(name).up.sql | sort | tail -n 1); \
	down_file="$${up_file%.up.sql}.down.sql"; \
	printf "%s\n" "$$diff_sql" > "$$up_file"; \
	if [ -n "$$reverse_sql" ]; then \
		printf "%s\n" "$$reverse_sql" > "$$down_file"; \
	fi; \
	echo "wrote $$up_file"; \
	echo "wrote $$down_file"

check-schema-drift:
	@set -e; \
	$(MAKE) atlas-reset-dev-db; \
	diff_sql="$$(ATLAS_DB_URL='$(atlasDbUrl)' ATLAS_DEV_URL='$(atlasDevUrl)' docker-compose --env-file .env -f $(composeFile) run --rm atlas schema diff --env local --from "$(atlasDbUrl)" --to env://src --exclude "public.schema_migrations" --format '{{ sql . }}')"; \
	if [ -n "$$diff_sql" ]; then \
		echo "schema drift detected between migrations/applied db and internal/db/persistence"; \
		echo "$$diff_sql"; \
		exit 1; \
	fi; \
	echo "no schema drift detected"

feature-verify: check-schema-drift verify-gen-query test
