module github.com/zcking/zerobus-examples/go-ingest-cli

go 1.25.6

require (
	github.com/databricks/zerobus-sdk-go v0.2.1 // indirect
	github.com/zcking/zerobus-examples/go-ingest-cli v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// Use local module
replace github.com/zcking/zerobus-examples/go-ingest-cli => ./
