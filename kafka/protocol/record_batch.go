//nolint
package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

import (
	"time"
)

const recordBatchOverhead = 49

type recordsArray []*Record

func (e recordsArray) decode(pd PacketDecoder) error {
	for i := range e {
		rec := &Record{}
		if err := rec.decode(pd); err != nil {
			return err
		}
		e[i] = rec
	}
	return nil
}

type RecordBatch struct {
	FirstOffset           int64
	PartitionLeaderEpoch  int32
	Version               int8
	Codec                 CompressionCodec
	CompressionLevel      int
	Control               bool
	LogAppendTime         bool
	LastOffsetDelta       int32
	FirstTimestamp        time.Time
	MaxTimestamp          time.Time
	ProducerID            int64
	ProducerEpoch         int16
	FirstSequence         int32
	Records               []*Record
	PartialTrailingRecord bool
	IsTransactional       bool

	compressedRecords []byte
	recordsLen        int // uncompressed Records size
}

func (b *RecordBatch) decode(pd PacketDecoder) (err error) {
	if b.FirstOffset, err = pd.getInt64(); err != nil {
		return err
	}

	batchLen, err := pd.getInt32()
	if err != nil {
		return err
	}

	if b.PartitionLeaderEpoch, err = pd.getInt32(); err != nil {
		return err
	}

	if b.Version, err = pd.getInt8(); err != nil {
		return err
	}

	crc32Decoder := acquireCrc32Field(crcCastagnoli)
	defer releaseCrc32Field(crc32Decoder)

	if err = pd.push(crc32Decoder); err != nil {
		return err
	}

	attributes, err := pd.getInt16()
	if err != nil {
		return err
	}
	b.Codec = CompressionCodec(int8(attributes) & compressionCodecMask)
	b.Control = attributes&controlMask == controlMask
	b.LogAppendTime = attributes&timestampTypeMask == timestampTypeMask
	b.IsTransactional = attributes&isTransactionalMask == isTransactionalMask

	if b.LastOffsetDelta, err = pd.getInt32(); err != nil {
		return err
	}

	if err = (Timestamp{&b.FirstTimestamp}).decode(pd); err != nil {
		return err
	}

	if err = (Timestamp{&b.MaxTimestamp}).decode(pd); err != nil {
		return err
	}

	if b.ProducerID, err = pd.getInt64(); err != nil {
		return err
	}

	if b.ProducerEpoch, err = pd.getInt16(); err != nil {
		return err
	}

	if b.FirstSequence, err = pd.getInt32(); err != nil {
		return err
	}

	numRecs, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if numRecs >= 0 {
		b.Records = make([]*Record, numRecs)
	}

	bufSize := int(batchLen) - recordBatchOverhead
	recBuffer, err := pd.getRawBytes(bufSize)
	if err != nil {
		if err == ErrInsufficientData {
			b.PartialTrailingRecord = true
			b.Records = nil
			return nil
		}
		return err
	}

	if err = pd.pop(); err != nil {
		return err
	}

	recBuffer, err = decompress(b.Codec, recBuffer)
	if err != nil {
		return err
	}

	b.recordsLen = len(recBuffer)
	err = decode(recBuffer, recordsArray(b.Records))
	if err == ErrInsufficientData {
		b.PartialTrailingRecord = true
		b.Records = nil
		return nil
	}
	return err
}
