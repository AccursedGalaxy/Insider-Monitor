package websocket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

// Binary message format:
// [1 byte magic number: 0x01]
// [1 byte version: 0x01]
// [1 byte message type]
// [4 bytes topic length]
// [N bytes topic]
// [8 bytes timestamp (int64)]
// [4 bytes payload length]
// [N bytes payload (JSON or compressed JSON)]
// [1 byte checksum (optional)]

const (
	// BinaryMagic is the magic number to identify binary messages
	BinaryMagic byte = 0x01
	// BinaryVersion is the current binary protocol version
	BinaryVersion byte = 0x01
)

// MessageEncoder defines interface for message encoding
type MessageEncoder interface {
	Encode(message *Message) ([]byte, error)
	Decode(data []byte) (*Message, error)
}

// JSONEncoder is the default JSON encoder
type JSONEncoder struct{}

// Encode encodes a message as JSON
func (e *JSONEncoder) Encode(message *Message) ([]byte, error) {
	return json.Marshal(message)
}

// Decode decodes a JSON message
func (e *JSONEncoder) Decode(data []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(data, &message)
	return &message, err
}

// BinaryEncoder encodes messages in a compact binary format
type BinaryEncoder struct {
	// Whether to compress the payload JSON
	Compress bool
	// Compression level to use
	CompressionLevel CompressionLevel
}

// Encode encodes a message in binary format
func (e *BinaryEncoder) Encode(message *Message) ([]byte, error) {
	var buf bytes.Buffer

	// Write magic number and version
	buf.WriteByte(BinaryMagic)
	buf.WriteByte(BinaryVersion)

	// Write message type (convert string to byte)
	typeID := messageTypeToID(string(message.Type))
	buf.WriteByte(typeID)

	// Write topic (write length first, then string)
	writeLengthPrefixedString(&buf, message.Topic)

	// Write timestamp (Unix timestamp in nanoseconds)
	timestamp := message.Time.UnixNano()
	var timeBytes [8]byte
	binary.BigEndian.PutUint64(timeBytes[:], uint64(timestamp))
	buf.Write(timeBytes[:])

	// Encode payload
	var payloadBytes []byte
	var err error

	// Convert payload to JSON bytes
	if payloadBytes, err = json.Marshal(message.Payload); err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	// Compress payload if requested
	if e.Compress && len(payloadBytes) > 256 {
		compressed, err := CompressData(payloadBytes, e.CompressionLevel)
		if err != nil {
			return nil, fmt.Errorf("failed to compress payload: %w", err)
		}

		// Only use compression if it saved space
		if len(compressed) < len(payloadBytes) {
			payloadBytes = compressed
		}
	}

	// Write payload length and data
	payloadLen := uint32(len(payloadBytes))
	var lenBytes [4]byte
	binary.BigEndian.PutUint32(lenBytes[:], payloadLen)
	buf.Write(lenBytes[:])
	buf.Write(payloadBytes)

	// Optionally add a simple checksum (sum of all previous bytes modulo 256)
	if false { // Disabled for now
		var checksum byte
		data := buf.Bytes()
		for _, b := range data {
			checksum += b
		}
		buf.WriteByte(checksum)
	}

	return buf.Bytes(), nil
}

// Decode decodes a binary message
func (e *BinaryEncoder) Decode(data []byte) (*Message, error) {
	if len(data) < 15 { // Minimum length for a valid message
		return nil, fmt.Errorf("message too short")
	}

	// Check magic number and version
	if data[0] != BinaryMagic || data[1] != BinaryVersion {
		return nil, fmt.Errorf("invalid binary message format")
	}

	// Read message type
	typeID := data[2]
	messageType := idToMessageType(typeID)

	// Read topic
	topicLen := int(binary.BigEndian.Uint32(data[3:7]))
	if 7+topicLen > len(data) {
		return nil, fmt.Errorf("invalid topic length")
	}
	topic := string(data[7 : 7+topicLen])

	// Read timestamp
	timeOffset := 7 + topicLen
	timestamp := int64(binary.BigEndian.Uint64(data[timeOffset : timeOffset+8]))
	messageTime := time.Unix(0, timestamp)

	// Read payload
	payloadLenOffset := timeOffset + 8
	payloadLen := int(binary.BigEndian.Uint32(data[payloadLenOffset : payloadLenOffset+4]))
	payloadOffset := payloadLenOffset + 4
	if payloadOffset+payloadLen > len(data) {
		return nil, fmt.Errorf("invalid payload length")
	}
	payloadData := data[payloadOffset : payloadOffset+payloadLen]

	// Check if payload is compressed (gzip magic bytes)
	var payload map[string]interface{}
	if len(payloadData) >= 2 && payloadData[0] == 0x1f && payloadData[1] == 0x8b {
		// Decompress payload
		decompressed, err := DecompressData(payloadData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress payload: %w", err)
		}
		// Parse decompressed JSON
		if err := json.Unmarshal(decompressed, &payload); err != nil {
			return nil, fmt.Errorf("failed to decode payload JSON: %w", err)
		}
	} else {
		// Parse regular JSON
		if err := json.Unmarshal(payloadData, &payload); err != nil {
			return nil, fmt.Errorf("failed to decode payload JSON: %w", err)
		}
	}

	// Create message
	message := &Message{
		Type:    MessageType(messageType),
		Topic:   topic,
		Payload: payload,
		Time:    messageTime,
	}

	return message, nil
}

// Helper to write a length-prefixed string
func writeLengthPrefixedString(buf *bytes.Buffer, s string) {
	strLen := uint32(len(s))
	var lenBytes [4]byte
	binary.BigEndian.PutUint32(lenBytes[:], strLen)
	buf.Write(lenBytes[:])
	buf.WriteString(s)
}

// Map message types to byte identifiers
func messageTypeToID(msgType string) byte {
	typeMap := map[string]byte{
		string(WalletUpdateMsg): 1,
		string(ConfigUpdateMsg): 2,
		string(AlertMsg):        3,
		string(StatusUpdateMsg): 4,
		string(SubscribeMsg):    5,
		string(UnsubscribeMsg):  6,
		string(PingMsg):         7,
	}

	if id, ok := typeMap[msgType]; ok {
		return id
	}
	return 0 // Unknown type
}

// Map byte identifiers to message types
func idToMessageType(id byte) string {
	typeMap := map[byte]string{
		1: string(WalletUpdateMsg),
		2: string(ConfigUpdateMsg),
		3: string(AlertMsg),
		4: string(StatusUpdateMsg),
		5: string(SubscribeMsg),
		6: string(UnsubscribeMsg),
		7: string(PingMsg),
	}

	if msgType, ok := typeMap[id]; ok {
		return msgType
	}
	return "unknown" // Unknown ID
}
