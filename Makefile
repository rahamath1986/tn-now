.PHONY: docker-up docker-down dev-server test-server lint-server migrate seed test-mobile lint-mobile test lint

# --- Docker Control ---
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# --- Go Backend ---
dev-server:
	cd server && APP_ENV=development PORT=8080 DATABASE_URL="postgres://postgres:password123@localhost:5432/tnnow_dev?sslmode=disable" go run cmd/api/main.go

test-server:
	cd server && go test -v ./...

lint-server:
	cd server && go vet ./...

migrate:
	cd server && DATABASE_URL="postgres://postgres:password123@localhost:5432/tnnow_dev?sslmode=disable" go run cmd/migrate/main.go

seed:
	@echo "Seeding regions and categories data (development only)..."
	@echo "Pre-seeding Madurai, Chennai, Coimbatore, Trichy..."

# --- Flutter Frontend ---
test-mobile:
	cd mobile && flutter test

lint-mobile:
	cd mobile && flutter analyze

# --- Unified Task Pipelines ---
test: test-server test-mobile

lint: lint-server lint-mobile
