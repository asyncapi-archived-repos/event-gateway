package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

import "fmt"

const (
	unknownRecords = iota
	legacyRecords
	defaultRecords

	magicOffset = 16
)

// Records implements a union type containing either a RecordBatch or a legacy MessageSet.
type Records struct {
	recordsType int
	MsgSet      *MessageSet
	RecordBatch *RecordBatch
}

func (r *Records) setTypeFromMagic(pd PacketDecoder) error {
	magic, err := magicValue(pd)
	if err != nil {
		return err
	}

	r.recordsType = defaultRecords
	if magic < 2 {
		r.recordsType = legacyRecords
	}

	return nil
}

func (r *Records) decode(pd PacketDecoder) error {
	if r.recordsType == unknownRecords {
		if err := r.setTypeFromMagic(pd); err != nil {
			return err
		}
	}

	switch r.recordsType {
	case legacyRecords:
		r.MsgSet = &MessageSet{}
		return r.MsgSet.decode(pd)
	case defaultRecords:
		r.RecordBatch = &RecordBatch{}
		return r.RecordBatch.decode(pd)
	}
	return fmt.Errorf("unknown Records type: %v", r.recordsType)
}

func magicValue(pd PacketDecoder) (int8, error) {
	return pd.peekInt8(magicOffset)
}
