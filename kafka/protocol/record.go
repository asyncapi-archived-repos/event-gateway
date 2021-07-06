//nolint
package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

import (
	"time"
)

const (
	isTransactionalMask = 0x10
	controlMask         = 0x20
)

// RecordHeader stores key and value for a record header
type RecordHeader struct {
	Key   []byte
	Value []byte
}

func (h *RecordHeader) decode(pd PacketDecoder) (err error) {
	if h.Key, err = pd.getVarintBytes(); err != nil {
		return err
	}

	if h.Value, err = pd.getVarintBytes(); err != nil {
		return err
	}
	return nil
}

// Record is kafka record type
type Record struct {
	Headers []*RecordHeader

	Attributes     int8
	TimestampDelta time.Duration
	OffsetDelta    int64
	Key            []byte
	Value          []byte
	length         varintLengthField
}

func (r *Record) decode(pd PacketDecoder) (err error) {
	if err = pd.push(&r.length); err != nil {
		return err
	}

	if r.Attributes, err = pd.getInt8(); err != nil {
		return err
	}

	timestamp, err := pd.getVarint()
	if err != nil {
		return err
	}
	r.TimestampDelta = time.Duration(timestamp) * time.Millisecond

	if r.OffsetDelta, err = pd.getVarint(); err != nil {
		return err
	}

	if r.Key, err = pd.getVarintBytes(); err != nil {
		return err
	}

	if r.Value, err = pd.getVarintBytes(); err != nil {
		return err
	}

	numHeaders, err := pd.getVarint()
	if err != nil {
		return err
	}

	if numHeaders >= 0 {
		r.Headers = make([]*RecordHeader, numHeaders)
	}
	for i := int64(0); i < numHeaders; i++ {
		hdr := new(RecordHeader)
		if err := hdr.decode(pd); err != nil {
			return err
		}
		r.Headers[i] = hdr
	}

	return pd.pop()
}
