.PHONY: build env frontend backend

build: frontend backend

frontend:
	cd frontend && npm run build
	cd frontend && mv out ../server/frontend

backend:
	cd server && go generate ./...
	cd server && go build -o ../app

env:
	@grep -v '^#' .env | sed 's/^/export /' > .env.sh
