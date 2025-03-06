package websocket

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// CompressionLevel defines how aggressively to compress messages
type CompressionLevel int

const (
	// CompressionNone disables compression
	CompressionNone CompressionLevel = iota
	// CompressionFast prioritizes speed over compression ratio
	CompressionFast
	// CompressionBest prioritizes compression ratio over speed
	CompressionBest
	// CompressionDefault uses a balanced approach
	CompressionDefault
)

// toGzipLevel converts our compression level to gzip constants
func (c CompressionLevel) toGzipLevel() int {
	switch c {
	case CompressionNone:
		return gzip.NoCompression
	case CompressionFast:
		return gzip.BestSpeed
	case CompressionBest:
		return gzip.BestCompression
	case CompressionDefault:
		return gzip.DefaultCompression
	default:
		return gzip.DefaultCompression
	}
}

// CompressData compresses a byte array using gzip
func CompressData(data []byte, level CompressionLevel) ([]byte, error) {
	// Skip compression for very small payloads
	if len(data) < 256 {
		return data, nil
	}

	// Create a buffer to store compressed data
	var buf bytes.Buffer

	// Create a gzip writer with the specified compression level
	gzWriter, err := gzip.NewWriterLevel(&buf, level.toGzipLevel())
	if err != nil {
		return nil, err
	}

	// Write data to the gzip writer
	if _, err := gzWriter.Write(data); err != nil {
		return nil, err
	}

	// Close the gzip writer to flush any remaining data
	if err := gzWriter.Close(); err != nil {
		return nil, err
	}

	// Return the compressed data
	return buf.Bytes(), nil
}

// DecompressData decompresses gzip compressed data
func DecompressData(data []byte) ([]byte, error) {
	// Check for gzip magic number (0x1f, 0x8b)
	if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
		// Not compressed, return as is
		return data, nil
	}

	// Create a reader for the compressed data
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	// Read all the decompressed data
	return io.ReadAll(gzReader)
}

// CompressibleMessage is an interface for messages that can be compressed
type CompressibleMessage interface {
	// Compress compresses the message
	Compress(level CompressionLevel) error
	// Decompress decompresses the message
	Decompress() error
	// IsCompressed checks if the message is compressed
	IsCompressed() bool
	// CompressedSize returns the size of the compressed message
	CompressedSize() int
	// OriginalSize returns the size of the original message
	OriginalSize() int
}

// WithCompression adds compression capabilities to Message
func (m *Message) WithCompression() *CompressedMessage {
	return &CompressedMessage{
		Message:         *m,
		compressed:      false,
		originalSize:    0,
		compressedSize:  0,
		compressionRate: 0,
	}
}

// CompressedMessage wraps a Message with compression capabilities
type CompressedMessage struct {
	Message         Message `json:"message"`
	compressed      bool
	originalSize    int
	compressedSize  int
	compressionRate float64
}

// Compress compresses the message payload
func (cm *CompressedMessage) Compress(level CompressionLevel) error {
	// Skip if already compressed or no compression requested
	if cm.compressed || level == CompressionNone {
		return nil
	}

	// Marshal the payload to JSON
	payloadBytes, err := marshalJSON(cm.Message.Payload)
	if err != nil {
		return err
	}

	// Record original size
	cm.originalSize = len(payloadBytes)

	// Compress the payload
	compressedBytes, err := CompressData(payloadBytes, level)
	if err != nil {
		return err
	}

	// Record compressed size
	cm.compressedSize = len(compressedBytes)

	// Only use compression if it actually saves space
	if cm.compressedSize < cm.originalSize {
		// Calculate compression rate
		cm.compressionRate = 1 - float64(cm.compressedSize)/float64(cm.originalSize)
		cm.compressed = true

		// Replace payload with compressed data
		cm.Message.Payload = map[string]interface{}{
			"_compressed": true,
			"_data":       compressedBytes,
		}
	}

	return nil
}

// Decompress decompresses the message payload
func (cm *CompressedMessage) Decompress() error {
	// Skip if not compressed
	if !cm.IsCompressed() {
		return nil
	}

	// Extract compressed data
	compressedData, ok := cm.Message.Payload["_data"].([]byte)
	if !ok {
		return fmt.Errorf("invalid compressed data format")
	}

	// Decompress the data
	decompressedBytes, err := DecompressData(compressedData)
	if err != nil {
		return err
	}

	// Unmarshal the decompressed data back to the payload
	var payload map[string]interface{}
	if err := unmarshalJSON(decompressedBytes, &payload); err != nil {
		return err
	}

	// Replace the payload with the decompressed data
	cm.Message.Payload = payload
	cm.compressed = false

	return nil
}

// IsCompressed checks if the message is compressed
func (cm *CompressedMessage) IsCompressed() bool {
	if compressed, ok := cm.Message.Payload["_compressed"].(bool); ok && compressed {
		return true
	}
	return false
}

// CompressedSize returns the size of the compressed message
func (cm *CompressedMessage) CompressedSize() int {
	return cm.compressedSize
}

// OriginalSize returns the size of the original message
func (cm *CompressedMessage) OriginalSize() int {
	return cm.originalSize
}

// CompressionRate returns the compression rate (0-1)
func (cm *CompressedMessage) CompressionRate() float64 {
	return cm.compressionRate
}

// Helper functions for JSON marshaling/unmarshaling
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
