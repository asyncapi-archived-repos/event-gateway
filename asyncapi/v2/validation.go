package v2

import (
	"encoding/json"
	"fmt"

	"github.com/asyncapi/event-gateway/asyncapi"
	"github.com/asyncapi/event-gateway/proxy"
	"github.com/xeipuuv/gojsonschema"
)

func FromDocJSONSchemaMessageValidator(doc asyncapi.Document) (proxy.MessageValidator, error) {
	channels := doc.ApplicationSubscribableChannels()
	messageSchemas := make(map[string]gojsonschema.JSONLoader)
	for _, c := range channels {
		for _, o := range c.Operations() {
			if !o.IsApplicationSubscribing() {
				continue
			}

			// Assuming there is only one message per operation as per Asyncapi 2.x.x.
			// See https://github.com/asyncapi/event-gateway/issues/10
			if len(o.Messages()) > 1 {
				return nil, fmt.Errorf("can not generate message validation for operation %s. Reason: the operation has more than one message and we can't correlate which one is it", o.ID())
			}

			if len(o.Messages()) == 0 {
				return nil, fmt.Errorf("can not generate message validation for operation %s. Reason:. Operation has no message. This is totally unexpected", o.ID())
			}

			// Assuming there is only one message per operation and one operation of a particular type per Channel.
			// See https://github.com/asyncapi/event-gateway/issues/10
			msg := o.Messages()[0]

			raw, err := json.Marshal(msg.Payload())
			if err != nil {
				return nil, fmt.Errorf("error marshaling message payload for generating json schema for validation. Operation: %s, Message: %s", o.ID(), msg.Name())
			}

			messageSchemas[c.ID()] = gojsonschema.NewBytesLoader(raw)
		}
	}

	idProvider := func(msg *proxy.Message) string {
		// messageSchemas map is indexed by Channel name, so we need to tell the validator.
		return msg.Context.Channel
	}

	return proxy.JSONSchemaMessageValidator(messageSchemas, idProvider)
}
