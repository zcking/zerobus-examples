package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	zerobus "github.com/databricks/zerobus-sdk-go"
	"github.com/google/uuid"
	"github.com/zcking/zerobus-examples/go-ingest-cli/gen/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
)

// IngestConfig holds all configuration for an ingestion run.
type IngestConfig struct {
	Endpoint     string
	Host         string
	ClientID     string
	ClientSecret string
	Table        string
	RequestID    string
	Message      string // Single inline message; ignored when stdin is piped.
}

func (c IngestConfig) validate() error {
	if c.Endpoint == "" || c.Host == "" || c.ClientID == "" || c.ClientSecret == "" || c.Table == "" {
		return fmt.Errorf("missing required config: set via flags or env vars (ZEROBUS_ENDPOINT, DATABRICKS_HOST, DATABRICKS_CLIENT_ID, DATABRICKS_CLIENT_SECRET, TABLE_NAME)")
	}
	return nil
}

// stdinIsPipe returns true when stdin is connected to a pipe (not a terminal).
func stdinIsPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func ingest(cfg IngestConfig) error {
	if err := cfg.validate(); err != nil {
		return err
	}

	if cfg.RequestID == "" {
		cfg.RequestID = uuid.New().String()
	}

	// Create SDK instance.
	sdk, err := zerobus.NewZerobusSdk(cfg.Endpoint, cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to create SDK: %w", err)
	}
	defer sdk.Free()

	// Build the protobuf descriptor from generated code.
	fileDescProto := protodesc.ToFileDescriptorProto(pb.File_cli_message_proto)
	descriptorBytes, err := proto.Marshal(fileDescProto.MessageType[0])
	if err != nil {
		return fmt.Errorf("failed to marshal descriptor: %w", err)
	}

	options := zerobus.DefaultStreamConfigurationOptions()

	// Create stream.
	stream, err := sdk.CreateStream(
		zerobus.TableProperties{
			TableName:       cfg.Table,
			DescriptorProto: descriptorBytes,
		},
		cfg.ClientID,
		cfg.ClientSecret,
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	if stdinIsPipe() {
		return ingestFromReader(stream, os.Stdin, cfg.RequestID)
	}

	if cfg.Message == "" {
		return fmt.Errorf("no input: provide --message or pipe data via stdin")
	}
	return ingestSingle(stream, cfg.Message, cfg.RequestID)
}

// ingestSingle sends a single message and waits for acknowledgment.
func ingestSingle(stream *zerobus.ZerobusStream, msg, requestID string) error {
	data, err := marshalRecord(msg, requestID)
	if err != nil {
		return err
	}

	offset, err := stream.IngestRecordOffset(data)
	if err != nil {
		return fmt.Errorf("failed to ingest record: %w", err)
	}
	log.Printf("Record ingested, offset: %d", offset)

	if err := stream.WaitForOffset(offset); err != nil {
		return fmt.Errorf("failed to wait for acknowledgment: %w", err)
	}
	log.Println("Acknowledged!")
	return nil
}

// ingestFromReader reads lines from r and ingests each as a record.
func ingestFromReader(stream *zerobus.ZerobusStream, r io.Reader, requestID string) error {
	scanner := bufio.NewScanner(r)
	var count int
	var lastOffset int64

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		data, err := marshalRecord(line, requestID)
		if err != nil {
			log.Printf("Failed to marshal record %d: %v", count, err)
			continue
		}

		offset, err := stream.IngestRecordOffset(data)
		if err != nil {
			return fmt.Errorf("failed to ingest record %d: %w", count, err)
		}
		lastOffset = offset
		count++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	if count == 0 {
		log.Println("No records to ingest.")
		return nil
	}

	log.Printf("%d records ingested, waiting for acknowledgment...", count)
	if err := stream.WaitForOffset(lastOffset); err != nil {
		return fmt.Errorf("failed to wait for acknowledgment: %w", err)
	}
	log.Println("All records acknowledged!")
	return nil
}

func marshalRecord(msg, requestID string) ([]byte, error) {
	record := &pb.Wrapper{
		RequestId: proto.String(requestID),
		Msg:       proto.String(msg),
		EventTime: proto.Int64(time.Now().UnixMicro()),
	}
	data, err := proto.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return data, nil
}
