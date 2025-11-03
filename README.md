# Zerobus Rust Examples

A collection of example applications demonstrating how to use the [Databricks Zerobus SDK for Rust](https://github.com/databricks/zerobus-sdk-rs).

## Overview

Zerobus is Databricks' streaming ingestion service that allows you to ingest data into Unity Catalog tables using Protocol Buffers over gRPC. This repository contains practical examples to help you get started.

For comprehensive documentation, see the official [Databricks Zerobus Ingest documentation](https://docs.databricks.com/aws/en/ingestion/lakeflow-connect/zerobus-ingest?language=Rust%C2%A0SDK).

## Prerequisites

- Rust 1.70 or later
- A Databricks workspace with Zerobus enabled (contact your Databricks account representative if needed)
- Service principal with OAuth credentials (client ID and client secret)
- A Unity Catalog table configured for Zerobus ingestion
- Protocol Buffer compiler (`protoc`) for generating schemas

## Setup

### 1. Clone this repository

```bash
git clone git@github.com:zcking/zerobus-rust-examples.git
cd zerobus-rust-examples
```

### 2. Configure environment variables

Copy the example environment file and fill in your actual values:

```bash
cp .env.example .env
```

Edit `.env` with your specific configuration and credentials.

Export the environment variables from `.env` in your shell:  

```bash
export $(grep -v '^#' .env | grep -v '^$' | xargs)
```

**Creating a service principal:**
1. In your Databricks workspace, go to **Settings** > **Identity and Access**
2. Create a new service principal and generate credentials
3. Grant the required permissions for your target table:
   ```sql
   GRANT USE CATALOG ON CATALOG <catalog> TO `<service-principal-uuid>`;
   GRANT USE SCHEMA ON SCHEMA <catalog.schema> TO `<service-principal-uuid>`;
   GRANT MODIFY, SELECT ON TABLE <catalog.schema.table> TO `<service-principal-uuid>`;
   ```

### 3. Generate Protocol Buffer schemas

To generate Protocol Buffer schemas from your Unity Catalog table, you'll use the `generate_files` tool from the Zerobus Rust SDK repository:

#### Clone the Zerobus Rust SDK repository

```bash
# Clone to a temporary directory (you only need this for the tools)
cd ..
git clone https://github.com/databricks/zerobus-sdk-rs.git
cd zerobus-sdk-rs/tools/generate_files
```

#### Generate the schema files

```bash
# Generate .proto file, Rust code, and descriptor from the Unity Catalog table
cargo run -- \
  --uc-endpoint $DATABRICKS_HOST \
  --client-id $DATABRICKS_CLIENT_ID \
  --client-secret $DATABRICKS_CLIENT_SECRET \
  --table $TABLE_NAME \
  --output-dir "../zerobus-rust-examples/examples/hello-world/proto"
```

This will generate three files in the proto directory:
- `<table_name>.proto` - Protocol Buffer schema definition
- `<table_name>.rs` - Rust code generated from the schema
- `<table_name>.descriptor` - Binary descriptor file for runtime use

**Note:** The proto directory structure is already set up in the examples. The `generate_files` tool will create all necessary artifacts that match your Unity Catalog table schema. I currently don't commit the `*.descriptor` files to Git because they are binary files.

## Examples

| Example | Description |
|---------|-------------|
| [Hello World](examples/hello-world/README.md) | Basic example demonstrating the fundamental workflow of the Zerobus SDK, including SDK initialization, stream creation, message encoding, record ingestion, and graceful shutdown. |

## Project Structure

```
zerobus-rust-examples/
├── Cargo.toml              # Workspace configuration
├── README.md               # This file
├── .env.example            # Example environment configuration
└── examples/
    └── hello-world/        # Basic hello world example
        ├── Cargo.toml
        ├── README.md
        ├── proto/          # Protocol Buffer definitions
        └── src/
            └── main.rs
```

## Key Concepts

### SDK Initialization

The SDK requires two endpoints:
- **Zerobus Endpoint**: The gRPC endpoint for streaming data ingestion (format: `https://<workspace_id>.zerobus.<region>.cloud.databricks.com`)
- **Databricks Host**: Your workspace URL used for Unity Catalog authentication and table metadata

### Authentication

The SDK handles OAuth 2.0 authentication automatically using service principal credentials. You only need to provide:
- **Client ID**: Service principal application ID
- **Client Secret**: Service principal secret

These credentials are used to obtain and refresh access tokens as needed. The SDK manages token lifecycle internally.

### Stream Lifecycle

1. **Create Stream**: Opens an authenticated bidirectional gRPC stream to the Zerobus service
2. **Ingest Records**: Send Protocol Buffer encoded data representing table rows
3. **Acknowledgments**: Each record returns a future that resolves when the service acknowledges durability
4. **Flush**: Force pending records to be transmitted (useful before shutdown)
5. **Close**: Gracefully shutdown the stream, ensuring all records are acknowledged

### Protocol Buffers

The Zerobus service uses Protocol Buffers for efficient data serialization. Here's the workflow:

1. **Generate Schema**: Use the `generate_files` tool from [zerobus-sdk-rs](https://github.com/databricks/zerobus-sdk-rs) to automatically generate Protocol Buffer definitions from your Unity Catalog table schema
2. **Include Generated Code**: The tool creates three files:
   - `.proto` - Protocol Buffer schema definition
   - `.rs` - Rust code with message structs
   - `.descriptor` - Binary descriptor for runtime type information
3. **Encode and Send**: Create instances of your message structs, encode them, and send via the stream

## Configuration Options

The SDK supports various configuration options via `StreamConfigurationOptions`:

| Option | Default | Description |
|--------|---------|-------------|
| `max_inflight_records` | 50000 | Maximum number of unacknowledged records |
| `recovery` | true | Enable automatic stream recovery on failures |
| `recovery_timeout_ms` | 15000 | Timeout for recovery operations (ms) |
| `recovery_backoff_ms` | 2000 | Delay between recovery attempts (ms) |
| `recovery_retries` | 3 | Maximum number of recovery attempts |
| `flush_timeout_ms` | 300000 | Timeout for flush operations (ms) |
| `server_lack_of_ack_timeout_ms` | 60000 | Server acknowledgment timeout (ms) |

## Troubleshooting

### Common Issues

**Authentication errors:**
- Verify your service principal credentials are correct
- Ensure the service principal has the required permissions on the target table
- Check that your Databricks host URL is correct

**Connection errors:**
- Verify your Zerobus endpoint format: `https://<workspace_id>.zerobus.<region>.cloud.databricks.com`
- Ensure network connectivity to your Databricks workspace
- Check that Zerobus is enabled for your workspace

**Schema errors:**
- Regenerate your Protocol Buffer files if your table schema has changed
- Ensure the descriptor file matches your current table schema
- Verify field types match between your data and the schema

## Resources

- [Databricks Zerobus Ingest Documentation](https://docs.databricks.com/aws/en/ingestion/lakeflow-connect/zerobus-ingest?language=Rust%C2%A0SDK)
- [Zerobus Rust SDK Repository](https://github.com/databricks/zerobus-sdk-rs)
- [Databricks Unity Catalog Documentation](https://docs.databricks.com/unity-catalog/)
- [Protocol Buffers Documentation](https://protobuf.dev/)

## License

Apache-2.0
