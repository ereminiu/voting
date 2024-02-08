build:
	docker compose build voting-app

run:
	docker compose up voting-app

test:
	go test -v ./...

migrate:
	go run cmd/migrate/main.go -mode=$(mode) -cmd=$(cmd)

migrate-force:
	go run cmd/migrate/main.go -mode=$(mode) -cmd=force -v=$(v)

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' down

swag:
	swag init -g cmd/main.go
