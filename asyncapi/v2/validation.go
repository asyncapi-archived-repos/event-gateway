package v2

import (
	"encoding/json"
	"fmt"
	"strings"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/asyncapi/event-gateway/asyncapi"
	"github.com/asyncapi/event-gateway/message"
	"github.com/xeipuuv/gojsonschema"
)

func FromDocJSONSchemaMessageValidator(doc asyncapi.Document) (message.Validator, error) {
	channels := doc.ApplicationSubscribableChannels()
	messageSchemas := make(map[string]gojsonschema.JSONLoader)
	for _, c := range channels {
		for _, o := range c.Operations() {
			if !o.IsApplicationSubscribing() {
				continue
			}

			if len(o.Messages()) == 0 {
				return nil, fmt.Errorf("can not generate message validation for operation %s. Reason:. Operation has no message. This is totally unexpected", o.ID())
			}

			var payload asyncapi.Schema
			var messageNames string
			if len(o.Messages()) > 1 {
				// Meaning message payload is a Schema containing several payloads as `oneOf`.
				// Generating back just one Schema adding all payloads to oneOf field.
				msgs := o.Messages()
				oneOfSchemas := make([]asyncapi.Schema, len(msgs))
				names := make([]string, len(msgs))
				for i, msg := range msgs {
					oneOfSchemas[i] = msg.Payload()
					names[i] = msg.Name()
				}
				payload = &Schema{OneOfField: oneOfSchemas}
				messageNames = strings.Join(names, ", ")
			} else {
				payload = o.Messages()[0].Payload()
				messageNames = o.Messages()[0].Name()
			}

			raw, err := json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("error marshaling message payload for generating json schema for validation. Operation: %s, Messages: %s", o.ID(), messageNames)
			}

			messageSchemas[c.ID()] = gojsonschema.NewBytesLoader(raw)
		}
	}

	idProvider := func(msg *watermillmessage.Message) string {
		// messageSchemas map is indexed by Channel name, so we need to tell the validator from where it should get that info.
		return msg.Metadata.Get(message.MetadataChannel)
	}

	return message.JSONSchemaMessageValidator(messageSchemas, idProvider)
}
