# Hello World Example

A basic example demonstrating the fundamental workflow of using the Databricks Zerobus SDK for Rust.

For comprehensive documentation, see the official [Databricks Zerobus Ingest documentation](https://docs.databricks.com/aws/en/ingestion/lakeflow-connect/zerobus-ingest?language=Rust%C2%A0SDK).

## What This Example Does

This example walks through the complete lifecycle of a Zerobus ingestion session:

1. **SDK Initialization**: Creates a `ZerobusSdk` instance with your workspace endpoints
2. **Stream Creation**: Opens an authenticated connection to a Unity Catalog table
3. **Message Encoding**: Encodes a simple message using Protocol Buffers
4. **Record Ingestion**: Sends the encoded message to Zerobus
5. **Acknowledgment**: Waits for confirmation that the message was received
6. **Graceful Shutdown**: Flushes pending records and closes the stream

## Prerequisites

- Rust 1.70 or later
- A Databricks workspace with Zerobus enabled (contact your Databricks account representative if needed)
- Service principal with OAuth credentials (client ID and client secret)
- A Unity Catalog table to ingest data into

## Setup

### 1. Create or Identify Your Unity Catalog Table

First, create or identify the target table in Unity Catalog. For example:

```sql
CREATE TABLE main.default.zerobus_hello_world (
    msg STRING,
    timestamp BIGINT
) USING DELTA;
```

### 2. Create a Service Principal and Grant Permissions

1. In your Databricks workspace, go to **Settings** > **Identity and Access**
2. Create a new service principal
3. Generate and save the client ID and client secret (save the secret securely - it's shown only once)
4. Copy the Application ID (UUID)
5. Grant the required permissions:

```sql
-- Replace <catalog>, <schema>, <table> with your actual names
-- Replace <service-principal-uuid> with your service principal's Application ID
GRANT USE CATALOG ON CATALOG main TO `<service-principal-uuid>`;
GRANT USE SCHEMA ON SCHEMA main.default TO `<service-principal-uuid>`;
GRANT MODIFY, SELECT ON TABLE main.default.zerobus_hello_world TO `<service-principal-uuid>`;
```

### 3. Configure Environment Variables

From the repository root, copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your specific configuration.

### 4. Generate Protocol Buffer Schema

See the root [README](../../README.md) for instructions on generating the required protobuf artifacts for examples.

### 5. Run the Example

```bash
cd examples/hello-world
cargo run
```

Expected output:
```
Zerobus Hello World Example
=============================

Initializing Zerobus SDK...
Creating stream to table: main.default.zerobus_hello_world
Stream created successfully!

Sending message: Hello, Zerobus!
Message sent, waiting for acknowledgment...
Message acknowledged successfully!
Stream flushed.

Stream closed. Hello World example complete!
```