.PHONY: swagger

swagger:
	swag init -g main.go \
		-d cmd/api,internal/book,internal/author,internal/health,internal/shared/response \
		-o docs
