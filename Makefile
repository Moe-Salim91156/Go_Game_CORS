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

# Run this from ~/Go_Game_CORS
test-e2e:
	@echo "Stopping containers..."
	docker-compose down
	@echo "Wiping persistent database..."
	rm -f backend/app.db
	@echo "Starting fresh backend..."
	docker-compose up -d
	@echo "Running Playwright..."
	npx playwright test --headed
