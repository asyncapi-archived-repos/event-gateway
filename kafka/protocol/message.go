//nolint
package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

import (
	"fmt"
	"time"
)

const (
	// CompressionNone no compression
	CompressionNone CompressionCodec = iota
	// CompressionGZIP compression using GZIP
	CompressionGZIP
	// CompressionSnappy compression using snappy
	CompressionSnappy
	// CompressionLZ4 compression using LZ4
	CompressionLZ4
	// CompressionZSTD compression using ZSTD
	CompressionZSTD

	// The lowest 3 bits contain the compression codec used for the message
	compressionCodecMask int8 = 0x07

	// Bit 3 set for "LogAppend" timestamps
	timestampTypeMask = 0x08
)

// CompressionCodec represents the various compression codecs recognized by Kafka in messages.
type CompressionCodec int8

func (cc CompressionCodec) String() string {
	return []string{
		"none",
		"gzip",
		"snappy",
		"lz4",
		"zstd",
	}[int(cc)]
}

// Message is a kafka message type
type Message struct {
	Codec            CompressionCodec // codec used to compress the message contents
	CompressionLevel int              // compression level
	LogAppendTime    bool             // the used timestamp is LogAppendTime
	Key              []byte           // the message key, may be nil
	Value            []byte           // the message contents
	Set              *MessageSet      // the message set a message might wrap
	Version          int8             // v1 requires Kafka 0.10
	Timestamp        time.Time        // the timestamp of the message (version 1+ only)

	compressedCache []byte
	compressedSize  int // used for computing the compression ratio metrics
}

func (m *Message) decode(pd PacketDecoder) (err error) {
	crc32Decoder := acquireCrc32Field(crcIEEE)
	defer releaseCrc32Field(crc32Decoder)

	err = pd.push(crc32Decoder)
	if err != nil {
		return err
	}

	m.Version, err = pd.getInt8()
	if err != nil {
		return err
	}

	if m.Version > 1 {
		return PacketDecodingError{fmt.Sprintf("unknown magic byte (%v)", m.Version)}
	}

	attribute, err := pd.getInt8()
	if err != nil {
		return err
	}
	m.Codec = CompressionCodec(attribute & compressionCodecMask)
	m.LogAppendTime = attribute&timestampTypeMask == timestampTypeMask

	if m.Version == 1 {
		if err := (Timestamp{&m.Timestamp}).decode(pd); err != nil {
			return err
		}
	}

	m.Key, err = pd.getBytes()
	if err != nil {
		return err
	}

	m.Value, err = pd.getBytes()
	if err != nil {
		return err
	}

	// Required for deep equal assertion during tests but might be useful
	// for future metrics about the compression ratio in fetch requests
	m.compressedSize = len(m.Value)

	if m.Value != nil && m.Codec != CompressionNone {
		m.Value, err = decompress(m.Codec, m.Value)
		if err != nil {
			return err
		}

		if err := m.decodeSet(); err != nil {
			return err
		}
	}

	return pd.pop()
}

// decodes a message set from a previously encoded bulk-message
func (m *Message) decodeSet() (err error) {
	pd := realDecoder{raw: m.Value}
	m.Set = &MessageSet{}
	return m.Set.decode(&pd)
}
