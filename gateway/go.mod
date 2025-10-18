module github.com/go-mockingcode/gateway

go 1.24.0

require (
	github.com/go-mockingcode/logger v0.0.0
	github.com/joho/godotenv v1.5.1
)

replace github.com/go-mockingcode/logger => ../pkg/logger
