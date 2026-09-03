dbUrl ?= $(or $(MIGRATE_DATABASE_URL),postgres://myuser:mypassword@localhost:5433/myapp?sslmode=disable)
composeFile = "./deploy/docker-compose.yml"
dockerFile = "deploy/Dockerfile"

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
