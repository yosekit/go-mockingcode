.PHONY: generate clean

# Generate Go code from proto files
generate:
	@echo "Generating Go code from proto files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		auth.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		project.proto
	@echo "Code generation complete!"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f *.pb.go
	@echo "Clean complete!"

# Install protoc plugins if needed
install-tools:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed!"

