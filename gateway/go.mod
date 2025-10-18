module github.com/go-mockingcode/gateway

go 1.24.0

require (
	github.com/go-mockingcode/logger v0.0.0
	github.com/go-mockingcode/proto v0.0.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.70.0
)

require (
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250127172529-29210b9bc287 // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)

replace github.com/go-mockingcode/logger => ../pkg/logger

replace github.com/go-mockingcode/proto => ../pkg/proto
