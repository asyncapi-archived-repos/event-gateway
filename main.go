package main

import (
	"bytes"
	"io"
	"strings"

	"github.com/asyncapi/event-gateway/kafka"
	server "github.com/grepplabs/kafka-proxy/cmd/kafka-proxy"
	"github.com/grepplabs/kafka-proxy/proxy"
	"github.com/grepplabs/kafka-proxy/proxy/protocol"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type config struct {
	Debug                        bool
	KafkaProxyBrokersMapping     pipeSeparatedValues `required:"true" split_words:"true"`
	KafkaProxyBrokersDialMapping pipeSeparatedValues `split_words:"true"`
	KafkaProxyExtraFlags         pipeSeparatedValues `split_words:"true"`
}

type pipeSeparatedValues struct {
	values []string
}

func (b *pipeSeparatedValues) Set(value string) error { //nolint:unparam
	b.values = strings.Split(value, "|")
	return nil
}

type requestKeyHandler struct{}

func (r *requestKeyHandler) Handle(requestKeyVersion *protocol.RequestKeyVersion, src io.Reader, ctx *proxy.RequestsLoopContext, bufferRead *bytes.Buffer) (shouldReply bool, err error) {
	if requestKeyVersion.ApiKey != kafka.RequestAPIKeyProduce {
		return true, nil
	}

	shouldReply, err = proxy.DefaultProduceKeyHandlerFunc(requestKeyVersion, src, ctx, bufferRead)
	if err != nil {
		return
	}

	msg := make([]byte, int64(requestKeyVersion.Length-int32(4+len(bufferRead.Bytes()))))
	if _, err = io.ReadFull(io.TeeReader(src, bufferRead), msg); err != nil {
		return
	}

	var req kafka.ProduceRequest
	if err = kafka.VersionedDecode(msg, &req, requestKeyVersion.ApiVersion); err != nil {
		logrus.Errorln(errors.Wrap(err, "error decoding ProduceRequest"))

		// Do not return an error but log it.
		return shouldReply, nil
	}

	for _, r := range req.Records {
		for _, s := range r {
			if s.RecordBatch != nil {
				for _, r := range s.RecordBatch.Records {
					if !isValid(r.Value) {
						logrus.Errorln("Message is not valid")
					} else {
						logrus.Debugln("Message is valid")
					}
				}
			}
			if s.MsgSet != nil {
				for _, mb := range s.MsgSet.Messages {
					if !isValid(mb.Msg.Value) {
						logrus.Errorln("Message is not valid")
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

func main() {
	var c config
	if err := envconfig.Process("eventgateway", &c); err != nil {
		logrus.Fatal(err)
	}

	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		_ = server.Server.Flags().Set("log-level", "debug")
	}

	// Yeah, not a good practice at all but I guess it's fine for now.
	proxy.ActualDefaultRequestHandler.RequestKeyHandlers.Set(kafka.RequestAPIKeyProduce, &requestKeyHandler{})

	for _, v := range c.KafkaProxyExtraFlags.values {
		f := strings.Split(v, "=")
		_ = server.Server.Flags().Set(f[0], f[1])
	}

	for _, v := range c.KafkaProxyBrokersMapping.values {
		_ = server.Server.Flags().Set("bootstrap-server-mapping", v)
	}

	for _, v := range c.KafkaProxyBrokersDialMapping.values {
		_ = server.Server.Flags().Set("dial-address-mapping", v)
	}

	if err := server.Server.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
