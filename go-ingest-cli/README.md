
# Go Ingest CLI

A CLI for ingesting messages into a Databricks Unity Catalog table through Zerobus. Send a single inline message or pipe lines from stdin.

Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) for the CLI framework.

## Setup

### Table

Create the target table and grant permissions to your service principal:

```sql
CREATE CATALOG IF NOT EXISTS zerobus_examples;
CREATE SCHEMA IF NOT EXISTS zerobus_examples.go_ingest_cli;

GRANT USE CATALOG ON CATALOG zerobus_examples TO `<service-principal-uuid>`;
GRANT USE SCHEMA, MODIFY, SELECT ON SCHEMA zerobus_examples.go_ingest_cli TO `<service-principal-uuid>`;

CREATE TABLE IF NOT EXISTS zerobus_examples.go_ingest_cli.zerobus_messages (
  request_id STRING,
  msg STRING,
  event_time TIMESTAMP
)
TBLPROPERTIES (
  'delta.enableDeletionVectors' = 'true',
  'delta.enableRowTracking' = 'false',
  'delta.parquet.compression.codec' = 'zstd'
)
COMMENT 'Generic messages ingested through Zerobus'
;
```

### Environment Variables

Export credentials and endpoints. These are picked up automatically by the CLI:

```bash
export DATABRICKS_HOST="https://myworkspace.cloud.databricks.com"
export DATABRICKS_CLIENT_ID="<client-id>"
export DATABRICKS_CLIENT_SECRET="<client-secret>"
export ZEROBUS_ENDPOINT="https://<workspace_id>.zerobus.<region>.cloud.databricks.com"
```

### Proto Generation

Generate the protobuf schema from the Unity Catalog table:

```bash
export TABLE_NAME=zerobus_examples.go_ingest_cli.zerobus_messages
make proto
```

## Usage

```bash
❯ ./go-ingest-cli --help 
Ingest records into a Databricks Unity Catalog table via Zerobus.

Send a single message inline:
  go-ingest-cli --table catalog.schema.table --message "hello"

Or pipe lines from stdin:
  cat messages.txt | go-ingest-cli --table catalog.schema.table
  echo "single line" | go-ingest-cli --table catalog.schema.table

Usage:
  go-ingest-cli [flags]

Flags:
      --client-id string       Databricks service principal client ID
      --client-secret string   Databricks service principal client secret
      --endpoint string        Zerobus gRPC endpoint
  -h, --help                   help for go-ingest-cli
      --host string            Databricks workspace URL
      --message string         Single message to ingest (omit to read from stdin)
      --request-id string      Request ID for ingested messages (default: random UUID)
      --table string           Unity Catalog table name
```

All parameters can be set via **flags** or **environment variables**. Flags take precedence.

| Flag | Env Var | Description |
|------|---------|-------------|
| `--endpoint` | `ZEROBUS_ENDPOINT` | Zerobus gRPC endpoint |
| `--host` | `DATABRICKS_HOST` | Databricks workspace URL |
| `--client-id` | `DATABRICKS_CLIENT_ID` | Service principal client ID |
| `--client-secret` | `DATABRICKS_CLIENT_SECRET` | Service principal secret |
| `--table` | `TABLE_NAME` | Unity Catalog table name |
| `--request-id` | | Request ID for messages (default: random UUID) |
| `--message` | | Single inline message (omit to read from stdin) |

### Single Message

```bash
# Build the CLI first
make build

# Direct usage
./go-ingest-cli --table zerobus_examples.go_ingest_cli.zerobus_messages --message "hello world"

# Usage via make
make invoke TABLE=zerobus_examples.go_ingest_cli.zerobus_messages MSG="hello world"
```

### Piped Input (stdin)

Each line of stdin is ingested as a separate record. Empty lines are skipped.

```bash
# Pipe from a file
cat messages.txt | ./go-ingest-cli --table zerobus_examples.go_ingest_cli.zerobus_messages

# Pipe from another command
kubectl logs my-pod | ./go-ingest-cli --table zerobus_examples.go_ingest_cli.zerobus_messages

# Single line via echo
echo "one-off event" | ./go-ingest-cli --table zerobus_examples.go_ingest_cli.zerobus_messages
```

### Build

```bash
make build
```

### Help

```bash
./go-ingest-cli --help
```
