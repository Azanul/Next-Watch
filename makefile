env:
	@grep -v '^#' .env | sed 's/^/export /' > .env.sh