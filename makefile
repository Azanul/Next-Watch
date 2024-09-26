.PHONY: build env frontend backend graphql-schema

build: frontend backend

frontend:
	cd frontend && npm install && npm run build
	cd frontend && mv out ../server/frontend

backend:
	cd server && go generate ./...
	cd server && go build -o ../app

env:
	@grep -v '^#' .env | sed 's/^/export /' > .env.sh

graphql-schema:
	cd server && go get github.com/99designs/gqlgen
	cd server && go run github.com/99designs/gqlgen generate

migrate-up:
	migrate -path db/migration -database DATABASE_URL up

migrate-down:
	migrate -path db/migration -database DATABASE_URL down