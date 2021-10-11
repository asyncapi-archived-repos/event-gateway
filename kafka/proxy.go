package kafka

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Shopify/sarama"
	watermillkafka "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/asyncapi/event-gateway/message"
	"github.com/asyncapi/event-gateway/proxy"
	server "github.com/grepplabs/kafka-proxy/cmd/kafka-proxy"
	kafkaproxy "github.com/grepplabs/kafka-proxy/proxy"
	kafkaprotocol "github.com/grepplabs/kafka-proxy/proxy/protocol"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Kafka request API Keys. See https://kafka.apache.org/protocol#protocol_api_keys.
const (
	// RequestAPIKeyProduce is the Kafka request API Key for the Produce Request.
	RequestAPIKeyProduce = 0
)

const (
	messagesChannelName = "kafka-produced-messages"
)

var defaultMarshaler = watermillkafka.DefaultMarshaler{}

// NewProxy creates a new Kafka Proxy based on a given configuration.
func NewProxy(c *ProxyConfig, r *watermillmessage.Router) (proxy.Proxy, error) {
	if c == nil {
		return nil, errors.New("config should be provided")
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Yeah, not a good practice at all but I guess it's fine for now.
	kafkaproxy.ActualDefaultRequestHandler.RequestKeyHandlers.Set(RequestAPIKeyProduce, NewProduceRequestHandler(r, c.MessageHandler, c.MessagePublisher, c.PublishToTopic))

	// Setting some defaults.
	_ = server.Server.Flags().Set("default-listener-ip", "0.0.0.0") // Binding to all local network interfaces. Needed for external calls.

	if c.BrokersMapping == nil {
		return nil, errors.New("Brokers mapping is required")
	}

	if c.Debug {
		_ = server.Server.Flags().Set("log-level", "debug")
	}

	if c.TLS != nil && c.TLS.Enable {
		_ = server.Server.Flags().Set("tls-enable", "true")
		_ = server.Server.Flags().Set("tls-insecure-skip-verify", fmt.Sprintf("%v", c.TLS.InsecureSkipVerify))
		_ = server.Server.Flags().Set("tls-client-cert-file", c.TLS.ClientCertFile)
		_ = server.Server.Flags().Set("tls-client-key-file", c.TLS.ClientKeyFile)
		_ = server.Server.Flags().Set("tls-ca-chain-cert-file", c.TLS.CAChainCertFile)
	}

	for _, v := range c.ExtraConfig {
		f := strings.Split(v, "=")
		_ = server.Server.Flags().Set(f[0], f[1])
	}

	for _, v := range c.BrokersMapping {
		_ = server.Server.Flags().Set("bootstrap-server-mapping", v)
	}

	for _, v := range c.DialAddressMapping {
		_ = server.Server.Flags().Set("dial-address-mapping", v)
	}

	return func(ctx context.Context) error {
		return server.Server.Execute()
	}, nil
}

// NewProduceRequestHandler creates a new request key handler for the Produce Request.
func NewProduceRequestHandler(r *watermillmessage.Router, handler watermillmessage.HandlerFunc, publisher watermillmessage.Publisher, publishToTopic string) kafkaproxy.KeyHandler {
	if handler == nil {
		return &produceRequestHandler{}
	}

	chanConfig := gochannel.Config{
		OutputChannelBuffer: 100, // TODO consider making this configurable
	}

	goChannelPubSub := gochannel.NewGoChannel(chanConfig, message.NewWatermillLogrusLogger(logrus.StandardLogger()))
	if publisher == nil {
		// This time we use a noPublisher handler, so converting the given handler.
		h := func(msg *watermillmessage.Message) error {
			_, err := handler(msg)
			return err
		}
		r.AddNoPublisherHandler("on-produce-request", messagesChannelName, goChannelPubSub, h)
	} else {
		r.AddHandler("on-produce-request", messagesChannelName, goChannelPubSub, publishToTopic, publisher, handler)
	}

	return &produceRequestHandler{
		publisher: goChannelPubSub,
	}
}

type produceRequestHandler struct {
	publisher watermillmessage.Publisher
}

func (h *produceRequestHandler) Handle(requestKeyVersion *kafkaprotocol.RequestKeyVersion, src io.Reader, ctx *kafkaproxy.RequestsLoopContext, bufferRead *bytes.Buffer) (shouldReply bool, err error) {
	if h.publisher == nil {
		logrus.Infoln("No message publisher is set. Skipping produceRequestHandler")
		return true, nil
	}

	if requestKeyVersion.ApiKey != RequestAPIKeyProduce {
		return true, nil
	}

	// TODO error handling should be responsibility of an error handler instead of being just logged.
	shouldReply, err = kafkaproxy.DefaultProduceKeyHandlerFunc(requestKeyVersion, src, ctx, bufferRead)
	if err != nil {
		return
	}

	msg := make([]byte, int64(requestKeyVersion.Length-int32(4+bufferRead.Len())))
	if _, err = io.ReadFull(io.TeeReader(src, bufferRead), msg); err != nil {
		return
	}

	// Hack for making compatible greplabs/kafka-proxy processor with Shopify/sarama ProduceRequest.
	// As both Transactional ID and ACKs has been read already by the processor, we fake them here because the Sarama decoder expects them to be present.
	// This information is not going to be used later on, as this is a read-only message.
	// transactional_id_size: 255, 255 | acks: 0, 1
	// TODO is there a way this info can be subtracted from kafka-proxy?
	msg = append([]byte{255, 255, 0, 1}, msg...)

	var req sarama.ProduceRequest
	if err = sarama.DoVersionedDecode(msg, &req, requestKeyVersion.ApiVersion); err != nil {
		logrus.WithError(err).Error("error decoding ProduceRequest")
		return shouldReply, nil
	}

	msgs, err := h.extractMessages(req)
	if err != nil {
		logrus.WithError(err).Error("error extracting messages")
		return shouldReply, nil
	}

	if len(msgs) == 0 {
		logrus.Error("Unexpected error: The produce request has no messages")
		return
	}

	if err := h.publisher.Publish(messagesChannelName, msgs...); err != nil {
		logrus.WithError(err).Error("error handling message")
		return shouldReply, nil
	}

	return shouldReply, nil
}

func (h *produceRequestHandler) extractMessages(req sarama.ProduceRequest) ([]*watermillmessage.Message, error) {
	var msgs []*watermillmessage.Message
	for topic, records := range req.Records {
		for partition, s := range records {
			if s.RecordBatch != nil {
				for _, r := range s.RecordBatch.Records {
					msg, err := defaultMarshaler.Unmarshal(&sarama.ConsumerMessage{
						Headers:   r.Headers,
						Key:       r.Key,
						Value:     r.Value,
						Topic:     topic,
						Partition: partition,
					})
					if err != nil {
						return nil, err
					}

					// Injecting the current Channel (kafka topic here) into the message Metadata (where Kafka headers are stored as well).
					msg.Metadata.Set(message.MetadataChannel, topic)
					msg.UUID = string(r.Key)
					msgs = append(msgs, msg)
				}
			}
			if s.MsgSet != nil {
				for _, mb := range s.MsgSet.Messages {
					msg, err := defaultMarshaler.Unmarshal(&sarama.ConsumerMessage{
						Key:       mb.Msg.Key,
						Value:     mb.Msg.Value,
						Timestamp: mb.Msg.Timestamp,
						Topic:     topic,
						Partition: partition,
						Offset:    mb.Offset,
					})
					if err != nil {
						return nil, err
					}

					// Injecting the current Channel (kafka topic here) into the message Metadata (where Kafka headers are stored as well).
					msg.Metadata.Set(message.MetadataChannel, topic)
					msg.UUID = string(mb.Msg.Key)
					msgs = append(msgs, msg)
				}
			}
		}
	}
	return msgs, nil
}
