.PHONY: up down rebuild shell logs tests ps

up:
	docker-compose up -d

down:
	docker-compose down

rebuild: 
	docker-compose down
	docker-compose up -d --build

logs: 
	docker-compose logs

ps : 
	docker ps
	docker-compose ps

shell-backend: ## Go inside the backend
	docker exec -it $$(docker ps -q -f name=backend) /bin/sh

shell-frontend: ## Go inside the frontend
	docker exec -it $$(docker ps -q -f name=frontend) /bin/sh

tests: 
	npx playwright test --headed

test-e2e:
	@echo "Stopping containers..."
	docker-compose down
	@echo "Wiping persistent database..."
	rm -f backend/app.db
	@echo "Starting fresh containers..."
	docker-compose up -d --build
	@echo "Waiting for backend to be ready..."
	@until curl -sf http://localhost:8080/api/health > /dev/null 2>&1; do \
		echo "  ...backend not ready yet, retrying in 1s"; \
		sleep 1; \
	done
	@echo "Backend is up!"
	@echo "Running Playwright..."
	npx playwright test --headed
