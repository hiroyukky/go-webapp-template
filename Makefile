
.PHONY: destroy
destroy:
	docker-compose down --rmi all --volumes --remove-orphans

.PHONY: build
build:
	docker-compose build --no-cache && docker-compose up

.PHONY: run
run:
	go run cmd/app/main.go

