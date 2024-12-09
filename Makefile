COMPOSE_FILE = docker-compose.yaml

.PHONY: build up down restart style

build: 
	docker-compose -f $(COMPOSE_FILE) build

up: 
	docker-compose -f $(COMPOSE_FILE) up -d

down: 
	docker-compose -f $(COMPOSE_FILE) down

restart: down build up

style:
	go fmt ./...