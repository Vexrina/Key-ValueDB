COMPOSE_FILE = docker-compose.yaml

.PHONY: up down restart style

up: 
	docker-compose -f $(COMPOSE_FILE) up -d --build

down: 
	docker-compose -f $(COMPOSE_FILE) down

restart: down up

style:
	go fmt ./...