# Default target
.PHONY: help
help:
	@echo "Generic CLI Message Ingestor - Available commands:"
	@echo ""
	@echo "Build & Package:"
	@echo "  make build           - Build the example"
	@echo "  make clean           - Clean build artifacts"
	@echo ""
	@echo "Protocol Buffers:"
	@echo "  make proto           - Generate proto files from Unity Catalog table"
	@echo "                        (requires DATABRICKS_HOST, DATABRICKS_CLIENT_ID,"
	@echo "                         DATABRICKS_CLIENT_SECRET, TABLE_NAME)"
	@echo ""
	@echo "Testing:"
	@echo "  make invoke MSG=hello - Invoke the CLI ingesting a simple 'hello' message"

# Variables
PROTO_DIR = proto
GEN_DIR = gen

# Build Lambda function
.PHONY: build
build: ARGS = --arm64
build:
	@echo "Building example..."
	@if ! command -v go &> /dev/null; then \
		echo "Error: go is not installed."; \
		exit 1; \
	fi
	go build -o . ./...
	@echo "Build complete!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning generated code..."
	rm -rf $(GEN_DIR)
	@echo "Clean complete!"

# Full proto workflow: generate .proto from UC, then compile with buf
.PHONY: proto
proto: proto-generate proto-compile

# Step 1: Generate .proto from Unity Catalog table using zerobus-generate
.PHONY: proto-generate
proto-generate:
	@echo "Generating Protocol Buffer files..."
	@if ! command -v zerobus-generate &> /dev/null; then \
		echo "Error: zerobus-generate is not installed."; \
		echo "Install it by:"; \
		echo "  1. Clone: git clone https://github.com/databricks/zerobus-sdk-rs.git"; \
		echo "  2. Build: cd zerobus-sdk-rs/rust/tools/generate_files && cargo build --release --target-dir=./target"; \
		echo "  3. Install: cp target/release/tools ~/.cargo/bin/zerobus-generate"; \
		exit 1; \
	fi
	@if [ -z "$$DATABRICKS_HOST" ] || [ -z "$$DATABRICKS_CLIENT_ID" ] || [ -z "$$DATABRICKS_CLIENT_SECRET" ] || [ -z "$$TABLE_NAME" ]; then \
		echo "Error: Required environment variables not set:"; \
		echo "  DATABRICKS_HOST"; \
		echo "  DATABRICKS_CLIENT_ID"; \
		echo "  DATABRICKS_CLIENT_SECRET"; \
		echo "  TABLE_NAME"; \
		exit 1; \
	fi
	zerobus-generate \
		--uc-endpoint $$DATABRICKS_HOST \
		--client-id $$DATABRICKS_CLIENT_ID \
		--client-secret $$DATABRICKS_CLIENT_SECRET \
		--table $$TABLE_NAME \
		--output-dir $(PROTO_DIR)
	@echo "Cleaning up old generated .rs and .descriptor files..."
	@rm -f $(PROTO_DIR)/*.go $(PROTO_DIR)/*.descriptor
	@echo "Proto files generated in $(PROTO_DIR)/"
	@echo "Note: Old .go and .descriptor files removed. Run 'make proto-compile' to regenerate with buf."

# Step 2: Compile .proto to language bindings and descriptor files using buf
.PHONY: proto-compile
proto-compile:
	@echo "Compiling proto files with buf..."
	@if ! command -v buf &> /dev/null; then \
		echo "Error: buf is not installed."; \
		echo "Install it with:"; \
		echo "  macOS: brew install bufbuild/buf/buf"; \
		echo "  Linux: https://buf.build/docs/installation"; \
		exit 1; \
	fi
	@echo "Generating language bindings..."
	buf generate $(PROTO_DIR)/
	@echo "Generating descriptor files..."
	@mkdir -p $(GEN_DIR)/descriptors
	@for proto_file in $(PROTO_DIR)/*.proto; do \
		if [ -f "$$proto_file" ]; then \
			base_name=$$(basename "$$proto_file" .proto); \
			buf build "$$proto_file" -o "$(GEN_DIR)/descriptors/$${base_name}.descriptor" --as-file-descriptor-set; \
		fi; \
	done
	@echo "Generated code in $(GEN_DIR)/"
	@echo "  - Descriptors: $(GEN_DIR)/descriptors/"

.PHONY: invoke
invoke: MSG ?= test
invoke: TABLE ?= zerobus_examples.go_ingest_cli.zerobus_messages
invoke:
	@echo "Ingesting message..."
	TABLE_NAME=$(TABLE) go run main.go

