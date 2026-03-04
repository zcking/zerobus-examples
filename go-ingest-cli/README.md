
# Go Ingest CLI

This example provides a generic, reusable, CLI for ingesting individual records, or streams of data, to a Databricks table through Zerobus.

## Setup

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

## Usage

For a simple test:  

1. `cd go-ingest-cli`
2. `make invoke TABLE=zerobus_examples.go_ingest_cli.zerobus_messages`

