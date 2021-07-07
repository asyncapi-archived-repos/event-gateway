package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

type ProduceRequest struct {
	Timeout int32
	Version int16 // v1 requires Kafka 0.9, v2 requires Kafka 0.10, v3 requires Kafka 0.11
	Records map[string]map[int32]Records
}

func (r *ProduceRequest) decode(pd PacketDecoder, version int16) error {
	r.Version = version

	// TransactionalID and Acks already read

	var err error
	if r.Timeout, err = pd.getInt32(); err != nil {
		return err
	}
	topicCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if topicCount == 0 {
		return nil
	}

	r.Records = make(map[string]map[int32]Records)
	for i := 0; i < topicCount; i++ {
		topic, err := pd.getString()
		if err != nil {
			return err
		}
		partitionCount, err := pd.getArrayLength()
		if err != nil {
			return err
		}
		r.Records[topic] = make(map[int32]Records)

		for j := 0; j < partitionCount; j++ {
			partition, err := pd.getInt32()
			if err != nil {
				return err
			}
			size, err := pd.getInt32()
			if err != nil {
				return err
			}
			recordsDecoder, err := pd.getSubset(int(size))
			if err != nil {
				return err
			}
			var records Records
			if err := records.decode(recordsDecoder); err != nil {
				return err
			}
			r.Records[topic][partition] = records
		}
	}

	return nil
}
