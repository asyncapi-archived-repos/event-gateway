package kafka

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/asyncapi/event-gateway/kafka/protocol"
	"github.com/asyncapi/event-gateway/proxy"
	server "github.com/grepplabs/kafka-proxy/cmd/kafka-proxy"
	kafkaproxy "github.com/grepplabs/kafka-proxy/proxy"
	kafkaprotocol "github.com/grepplabs/kafka-proxy/proxy/protocol"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewProxy creates a new Kafka Proxy based on a given configuration.
func NewProxy(c *ProxyConfig) (proxy.Proxy, error) {
	if c == nil {
		return nil, errors.New("config should be provided")
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Yeah, not a good practice at all but I guess it's fine for now.
	kafkaproxy.ActualDefaultRequestHandler.RequestKeyHandlers.Set(protocol.RequestAPIKeyProduce, &requestKeyHandler{})

	if c.BrokersMapping == nil {
		return nil, errors.New("Brokers mapping is required")
	}

	if c.Debug {
		_ = server.Server.Flags().Set("log-level", "debug")
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

	return func(_ context.Context) error {
		return server.Server.Execute()
	}, nil
}

type requestKeyHandler struct{}

func (r *requestKeyHandler) Handle(requestKeyVersion *kafkaprotocol.RequestKeyVersion, src io.Reader, ctx *kafkaproxy.RequestsLoopContext, bufferRead *bytes.Buffer) (shouldReply bool, err error) {
	if requestKeyVersion.ApiKey != protocol.RequestAPIKeyProduce {
		return true, nil
	}

	shouldReply, err = kafkaproxy.DefaultProduceKeyHandlerFunc(requestKeyVersion, src, ctx, bufferRead)
	if err != nil {
		return
	}

	msg := make([]byte, int64(requestKeyVersion.Length-int32(4+bufferRead.Len())))
	if _, err = io.ReadFull(io.TeeReader(src, bufferRead), msg); err != nil {
		return
	}
	var req protocol.ProduceRequest
	if err = protocol.VersionedDecode(msg, &req, requestKeyVersion.ApiVersion); err != nil {
		logrus.Errorln(errors.Wrap(err, "error decoding ProduceRequest"))
		// TODO notify error to a given notifier

		// Do not return an error but log it.
		return shouldReply, nil
	}

	for _, r := range req.Records {
		for _, s := range r {
			if s.RecordBatch != nil {
				for _, r := range s.RecordBatch.Records {
					if !isValid(r.Value) {
						logrus.Debugln("Message is not valid")
					} else {
						logrus.Debugln("Message is valid")
					}
				}
			}
			if s.MsgSet != nil {
				for _, mb := range s.MsgSet.Messages {
					if !isValid(mb.Msg.Value) {
						logrus.Debugln("Message is not valid")
					} else {
						logrus.Debugln("Message is valid")
					}
				}
			}
		}
	}

	return shouldReply, nil
}

func isValid(msg []byte) bool {
	return string(msg) != "invalid message"
}
