.PHONY: build env frontend backend

build:
	frontend
	mv out ../server/frontend
	backend

frontend:
	cd frontend && npm run build

backend:
	cd server && go generate ./...
	cd server && go build -o ../app

env:
	@grep -v '^#' .env | sed 's/^/export /' > .env.sh
