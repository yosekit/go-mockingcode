module github.com/go-mockingcode/data

go 1.24.0

require (
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/go-mockingcode/logger v0.0.0
	github.com/go-mockingcode/models v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.8.1
	go.mongodb.org/mongo-driver/v2 v2.3.1
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/go-mockingcode/logger => ../pkg/logger

replace github.com/go-mockingcode/models => ../pkg/models
