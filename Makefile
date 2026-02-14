.PHONY: up down rebuild shell seed

up:
	docker-compose up -d

down:
	docker-compose down

rebuild: 
	docker-compose down
	docker-compose up -d --build

shell-backend: ## Go inside the backend
	docker exec -it $$(docker ps -q -f name=backend) /bin/sh

shell-frontend: ## Go inside the frontend
	docker exec -it $$(docker ps -q -f name=frontend) /bin/sh

seed:
	cd backend/
	chmod +x seed.sh
	./seed.sh
